// Package web serves the single-file session viewer on loopback and bridges
// it to internal/sessions/live: the browser gets whole events over SSE and
// writes back through POST /send, while credentials and transport stay in
// this process.
//
// web.go is the entry point (Serve), server.go the HTTP routes and their
// guards, and hub.go the fan-out of live updates to connected pages.
// viewer/README.md says where the bundle and its wire protocol come from.
package web

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"

	"github.com/anthropics/anthropic-cli/v2/internal/sessions/live"
)

//go:embed viewer/viewer.html.gz
var viewerGzip []byte

// Session is what the server needs from an attached session; *live.Threads is one.
type Session interface {
	Updates() <-chan live.ThreadUpdate
	Snapshot() (threads []live.Thread, logs []live.ThreadLog)
	Events(id live.ThreadID) []live.Event
	Session() *anthropic.BetaManagedAgentsSession
	SendMessage(ctx context.Context, text string) error
	Interrupt(ctx context.Context) error
	Confirm(ctx context.Context, toolUseID string, allow bool, denyMessage string) error
}

var _ Session = (*live.Threads)(nil)

// Options configures Serve.
type Options struct {
	// ConsoleURL is the session's page in the web Console, for the header link.
	ConsoleURL string
	// Open is called once the server is listening, with the path of a local
	// page that leads to the viewer; nil just prints the viewer's URL via Logf.
	Open func(path string) error
	// Logf reports the viewer URL and browser-open failures to the user; required.
	Logf func(format string, args ...any)
}

// Serve serves the viewer for session on a loopback port until ctx is done.
func Serve(ctx context.Context, session Session, opts Options) error {
	page, err := loadViewer()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	srv := newServer(session, page, listener.Addr().String(), opts.ConsoleURL)
	go srv.hub.run()

	// The URL carries a one-time boot code, not the session token: it lands
	// in shell/browser history, so it must be worthless once used.
	url := fmt.Sprintf("http://%s/?boot=%s", listener.Addr(), srv.newBootCode())
	opts.Logf("Session viewer: %s", url)
	if opts.Open != nil {
		if err := srv.openFromFile(url, opts.Open); err != nil {
			opts.Logf("could not open a browser (%v); open the URL above", err)
		}
	}

	httpServer := &http.Server{Handler: srv, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()
	if err := httpServer.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

const redirectPage = `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Session viewer</title><meta http-equiv="refresh" content="0;url=%s"></head></html>`

// openFromFile opens url by way of a file that only this user can read: other
// local users can read the opener's arguments, and could spend the code first.
func (srv *server) openFromFile(url string, open func(path string) error) error {
	file, err := os.CreateTemp("", "ant-session-viewer-*.html")
	if err != nil {
		return err
	}
	srv.redirectFile = file.Name()
	_, err = fmt.Fprintf(file, redirectPage, url)
	if err = errors.Join(err, file.Close()); err == nil {
		err = open(file.Name())
	}
	if err != nil {
		srv.removeRedirectFile()
	}
	return err
}

func (srv *server) removeRedirectFile() {
	_ = os.Remove(srv.redirectFile)
}

// loadViewer inflates the embedded viewer page.
func loadViewer() ([]byte, error) {
	unzipped, err := gzip.NewReader(bytes.NewReader(viewerGzip))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(unzipped)
}
