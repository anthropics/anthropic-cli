package httpapi

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/nonce"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

// The fetch container of a sandbox pod logs why it was answered 403. It prints
// the answer only when the answer is a reason its script lists, so that the
// body of a 403 cannot put the pod's credential on that line. A reason the
// script does not list is logged as unknown, and the operator learns nothing
// from it. These tests hold the script's list to what this package writes.

const fetchScript = "../../../image/fetch.sh"

// listedReasons returns the reasons that refusal_reason of fetchScript lists,
// sorted.
func listedReasons(t *testing.T) []string {
	t.Helper()
	script, err := os.ReadFile(fetchScript)
	if err != nil {
		t.Fatal(err)
	}
	_, body, found := strings.Cut(string(script), "\nrefusal_reason() {\n")
	body, _, closed := strings.Cut(body, "\n}\n")
	if !found || !closed {
		t.Fatalf("%s has no refusal_reason function", fetchScript)
	}
	var listed []string
	for _, m := range regexp.MustCompile(`(?m)^\s*'([^']+)'( \| \\|\))$`).FindAllStringSubmatch(body, -1) {
		listed = append(listed, m[1])
	}
	if len(listed) == 0 {
		t.Fatalf("refusal_reason of %s lists no reason, one to a line in single quotes", fetchScript)
	}
	slices.Sort(listed)
	return listed
}

// Every way an Authenticator refuses a caller for good is played against the
// route. The script reads the first line of the answer, so each answer has to
// be its reason and one newline.
func TestTheFetchScriptListsTheReasonsOfA403AndNoOther(t *testing.T) {
	const theNonce, podIP = "the-nonce", "10.0.0.5"
	byNonce := func(locate nonce.PodLocatorFunc) Authenticator {
		registry := nonce.New(locate)
		registry.Set("sandbox-1", theNonce)
		return NonceAuth{Registry: registry}
	}
	located := func(running bool, err error) nonce.PodLocatorFunc {
		return func(context.Context, string) (types.UID, string, bool, error) { return "pod-a", podIP, running, err }
	}
	byToken := TokenAuth{Reviewer: fakeReviewer{"other-account": ErrNotSandbox, "no-pod": ErrNoPod}}
	written := map[string]bool{}
	for name, tc := range map[string]struct {
		auth              Authenticator
		bearer, from      string
		wantOtherThanA403 int
	}{
		"another account's token":        {auth: byToken, bearer: "other-account"},
		"a token bound to no pod":        {auth: byToken, bearer: "no-pod"},
		"a pod not seen running":         {auth: byNonce(located(false, nil)), from: podIP + ":5000"},
		"a pod that cannot be looked up": {auth: byNonce(located(true, errors.New("the watch is behind"))), from: podIP + ":5000"},
		"another address":                {auth: byNonce(located(true, nil)), from: "10.0.0.6:5000"},
		"an address with no port":        {auth: byNonce(located(true, nil)), from: podIP},
		// The pod's own request shows that the fixture refuses for the reason named.
		"the pod itself": {auth: byNonce(located(true, nil)), from: podIP + ":5000", wantOtherThanA403: http.StatusNotFound},
	} {
		r := requestFrom(t, theNonce, tc.from)
		if tc.bearer != "" {
			r.Header.Set("Authorization", "Bearer "+tc.bearer)
		}
		w := httptest.NewRecorder()
		(&Server{Auth: tc.auth, WorkItems: workitem.NewStore(), Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}).Handler().ServeHTTP(w, r)
		if tc.wantOtherThanA403 != 0 {
			if w.Code != tc.wantOtherThanA403 {
				t.Errorf("%s: status %d, want %d", name, w.Code, tc.wantOtherThanA403)
			}
			continue
		}
		reason, isLine := strings.CutSuffix(w.Body.String(), "\n")
		if w.Code != http.StatusForbidden || !isLine || strings.ContainsAny(reason, "\r\n") {
			t.Errorf("%s: the answer is %d %q, want 403 with one line that ends in a newline", name, w.Code, w.Body.String())
			continue
		}
		written[reason] = true
	}
	if got, listed := slices.Sorted(maps.Keys(written)), listedReasons(t); !slices.Equal(got, listed) {
		t.Errorf("a 403 carries one of the reasons\n%q\nand refusal_reason of %s lists\n%q\nGive the script the same list", got, fetchScript, listed)
	}
}

// notInA403 holds the error texts of this package and of package nonce that no
// 403 carries, each with where it ends up.
var notInA403 = map[string]string{
	"token is not valid for the dispatcher":                       "a 401",
	"a bearer token is required":                                  "a 401",
	"an assignment nonce is required":                             "a 401",
	"assignment nonce is not valid":                               "a 401",
	"token review: %w":                                            "the log, behind a 500 with a text of its own",
	"sandbox service account must be a dedicated account, not %q": "the start of the dispatcher",
	"token audience is required":                                  "the start of the dispatcher",
}

// The test above plays the refusals it knows. An error text added to either
// package could reach a 403 by a path that test does not play. So each text in
// their source is part of a reason the script lists, or is in notInA403.
func TestEveryErrorTextIsPartOfAListedReasonOrReachesNoA403(t *testing.T) {
	listed := listedReasons(t)
	inAListedReason := func(text string) bool {
		return slices.ContainsFunc(listed, func(reason string) bool {
			return reason == text || strings.HasPrefix(reason, text+": ") || strings.HasSuffix(reason, ": "+text)
		})
	}
	seen := map[string]bool{}
	sources, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	ofNonce, err := filepath.Glob("../nonce/*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range slices.Concat(sources, ofNonce) {
		if strings.HasSuffix(source, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), source, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			text, ok := errorText(n)
			// "%w: %v" words nothing itself. It joins two texts found here.
			if text = strings.TrimPrefix(text, "%w: "); !ok || text == "%v" {
				return true
			}
			seen[text] = true
			if _, elsewhere := notInA403[text]; !elsewhere && !inAListedReason(text) {
				t.Errorf("%s words the error %q. If a 403 can carry it, add the whole reason to refusal_reason of %s and play it in the test above. If none can, add it to notInA403", source, text, fetchScript)
			}
			return true
		})
	}
	for text := range notInA403 {
		if !seen[text] {
			t.Errorf("notInA403 holds %q, which neither package words", text)
		}
	}
	if len(seen) <= len(notInA403) {
		t.Errorf("the source words %d error texts, which leaves none for a 403: %v", len(seen), slices.Sorted(maps.Keys(seen)))
	}
}

// errorText returns the text n words, when n is a call of errors.New or
// fmt.Errorf, or a call of http.Error that answers 403, with a literal text.
func errorText(n ast.Node) (string, bool) {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return "", false
	}
	name := func(e ast.Expr) string {
		if sel, ok := e.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok {
				return pkg.Name + "." + sel.Sel.Name
			}
		}
		return ""
	}
	at := 0
	switch name(call.Fun) {
	case "errors.New", "fmt.Errorf":
	case "http.Error":
		if len(call.Args) != 3 || name(call.Args[2]) != "http.StatusForbidden" {
			return "", false
		}
		at = 1
	default:
		return "", false
	}
	lit, ok := call.Args[at].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	text, err := strconv.Unquote(lit.Value)
	return text, err == nil
}
