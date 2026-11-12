package web

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// server is the viewer's HTTP handler. GET / trades a one-time boot code for
// the page (which carries the session token); every other route wants that token.
type server struct {
	session    Session
	page       []byte
	host       string // "127.0.0.1:port"
	token      string
	consoleURL string
	hub        *hub

	bootMu    sync.Mutex
	bootCodes map[string]time.Time // unspent code → expiry

	redirectFile string // carries the first boot code to the browser; "" with no browser to open
}

func newServer(session Session, page []byte, host, consoleURL string) *server {
	return &server{
		session:    session,
		page:       page,
		host:       host,
		token:      randomToken(),
		consoleURL: consoleURL,
		hub:        newHub(session),
		bootCodes:  map[string]time.Time{},
	}
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !srv.hostAndOriginOK(r) {
		http.Error(w, "bad host", http.StatusForbidden)
		return
	}
	switch r.URL.Path {
	case "/":
		srv.serveIndex(w, r)
	case "/events":
		srv.serveEvents(w, r)
	case "/send":
		srv.serveSend(w, r)
	case "/reboot":
		srv.serveReboot(w, r)
	default:
		http.NotFound(w, r)
	}
}

// hostAndOriginOK rejects DNS-rebinding (Host must be our loopback address)
// and cross-site requests (Origin, when the browser sends one, must be us).
func (srv *server) hostAndOriginOK(r *http.Request) bool {
	_, port, _ := net.SplitHostPort(srv.host)
	if r.Host != srv.host && r.Host != "localhost:"+port {
		return false
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		return origin == "http://"+srv.host || origin == "http://localhost:"+port
	}
	return true
}

const (
	// nonceMarker is replaced with the per-response CSP nonce in every page we serve.
	nonceMarker = "__CSP_NONCE__"
	// configMarker is where the bundle expects its config <script>.
	configMarker = "<!--SESSION_VIEWER_CONFIG-->"

	csp = "default-src 'none'; script-src 'nonce-%[1]s'; style-src 'nonce-%[1]s'; font-src data:; img-src data: blob:; frame-src blob:; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'"

	// stashScript runs ahead of the viewer's own bootstrap on a freshly booted
	// page: it parks the session token in sessionStorage (per-origin, so
	// per-port — unlike a cookie, nothing else on 127.0.0.1 is ever sent it)
	// and scrubs the spent boot code from the address bar and history entry.
	stashScript = `<script nonce="__CSP_NONCE__">(function(){try{var c=JSON.parse(document.getElementById("session-viewer-config").textContent);sessionStorage.setItem("session-viewer-token",c.token)}catch(e){}history.replaceState(null,"","/")})()</script>`

	// resumePage is what a reload (no boot code) gets: if this tab still holds
	// the token it trades it for a fresh boot code and re-enters; otherwise the
	// link is dead and says so.
	resumePage = `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Session viewer</title><style nonce="__CSP_NONCE__">body{font:14px system-ui;margin:2rem}</style></head><body><p id="m">Reconnecting…</p><script nonce="__CSP_NONCE__">(function(){var m=document.getElementById("m"),t=sessionStorage.getItem("session-viewer-token");function dead(){m.textContent="This viewer link has expired. Run the ant command again for a new one."}if(!t){dead();return}fetch("/reboot",{method:"POST",headers:{Authorization:"Bearer "+t}}).then(function(r){return r.ok?r.json():Promise.reject()}).then(function(j){location.replace("/?boot="+encodeURIComponent(j.boot))}).catch(dead)})()</script></body></html>`
)

func (srv *server) serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	nonce := randomToken()
	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Content-Security-Policy", fmt.Sprintf(csp, nonce))
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Cache-Control", "no-store")

	if !srv.redeemBootCode(r.URL.Query().Get("boot")) {
		_, _ = w.Write(bytes.ReplaceAll([]byte(resumePage), []byte(nonceMarker), []byte(nonce)))
		return
	}
	srv.removeRedirectFile()
	// json.Marshal escapes <, > and & so the config cannot close its <script>.
	config, _ := json.Marshal(map[string]string{
		"eventsURL": "/events", "sendURL": "/send", "token": srv.token, "consoleURL": srv.consoleURL,
	})
	configScript := append(append([]byte(`<script type="application/json" id="session-viewer-config">`), config...), []byte("</script>"+stashScript)...)
	page := bytes.Replace(srv.page, []byte(configMarker), configScript, 1)
	page = bytes.ReplaceAll(page, []byte(nonceMarker), []byte(nonce))
	_, _ = w.Write(page)
}

// serveEvents streams the session as SSE: a snapshot, the last connection
// state, then every frame the hub broadcasts. EventSource cannot set headers,
// so the token rides in the query string.
func (srv *server) serveEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || !sameOrigin(r) || !srv.tokenOK(r.URL.Query().Get("token")) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-store")
	header.Set("X-Accel-Buffering", "no")

	// Subscribe before snapshotting so nothing falls between the two; a
	// frame that is also in the snapshot is a harmless upsert-by-id.
	frames := srv.hub.subscribe()
	defer srv.hub.unsubscribe(frames)
	write := func(f frame) error {
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", f.name, f.data); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}
	if write(srv.hub.snapshotFrame()) != nil || write(srv.hub.lastConnFrame()) != nil {
		return
	}
	for {
		select {
		case f, open := <-frames:
			if !open || write(f) != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

const maxSendBody = 64 << 10

// serveSend turns one user event from the page into the matching Session call.
func (srv *server) serveSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !sameOrigin(r) || !srv.tokenOK(bearerToken(r)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "send JSON", http.StatusUnsupportedMediaType)
		return
	}
	var body struct {
		Type        string `json:"type"`
		Text        string `json:"text"`
		ToolUseID   string `json:"tool_use_id"`
		Result      string `json:"result"`
		DenyMessage string `json:"deny_message"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, maxSendBody)).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	var err error
	switch body.Type {
	case "user.message":
		if strings.TrimSpace(body.Text) == "" {
			http.Error(w, "empty message", http.StatusBadRequest)
			return
		}
		err = srv.session.SendMessage(r.Context(), body.Text)
	case "user.interrupt":
		err = srv.session.Interrupt(r.Context())
	case "user.tool_confirmation":
		if body.ToolUseID == "" || (body.Result != "allow" && body.Result != "deny") {
			http.Error(w, "tool_use_id and result (allow|deny) required", http.StatusBadRequest)
			return
		}
		err = srv.session.Confirm(r.Context(), body.ToolUseID, body.Result == "allow", body.DenyMessage)
	default:
		http.Error(w, "unknown type "+body.Type, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	// A body (rather than 204) keeps browsers from logging the fetch as cancelled.
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

// serveReboot lets a tab that already holds the token get back in after a
// reload without the token ever riding in a URL.
func (srv *server) serveReboot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !sameOrigin(r) || !srv.tokenOK(bearerToken(r)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = fmt.Fprintf(w, `{"boot":%q}`, srv.newBootCode())
}

func (srv *server) tokenOK(got string) bool {
	return subtle.ConstantTimeCompare([]byte(got), []byte(srv.token)) == 1
}

func bearerToken(r *http.Request) string {
	token, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	return token
}

// sameOrigin holds the API routes to requests the page itself made. Browsers
// stamp every fetch/EventSource with Sec-Fetch-Site; a request without it, or
// from another site, is something else on this machine replaying a credential.
func sameOrigin(r *http.Request) bool {
	return r.Header.Get("Sec-Fetch-Site") == "same-origin"
}

const bootCodeTTL = 2 * time.Minute

// newBootCode mints a single-use code that GET / exchanges for the page (and,
// inside it, the session token).
func (srv *server) newBootCode() string {
	code := randomToken()
	srv.bootMu.Lock()
	defer srv.bootMu.Unlock()
	now := time.Now()
	for unspent, expiry := range srv.bootCodes {
		if now.After(expiry) {
			delete(srv.bootCodes, unspent)
		}
	}
	srv.bootCodes[code] = now.Add(bootCodeTTL)
	return code
}

// redeemBootCode spends a code: true at most once per code, and only in time.
func (srv *server) redeemBootCode(code string) bool {
	srv.bootMu.Lock()
	defer srv.bootMu.Unlock()
	expiry, ok := srv.bootCodes[code]
	delete(srv.bootCodes, code)
	return ok && time.Now().Before(expiry)
}

// randomToken returns 256 random bits, URL-safe.
func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
