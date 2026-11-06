package main

import "testing"

// The check matches net/http's and golang.org/x/net/http2's reading of
// GODEBUG. The warning then appears exactly when the HTTP/2 debug logging
// is on.
func TestHTTP2DebugOn(t *testing.T) {
	cases := []struct {
		godebug string
		want    bool
	}{
		{"", false},
		{"http2debug=0", false},
		{"http2debug=1", true},
		{"http2debug=2", true},
		{"gctrace=1,http2debug=1", true},
		{"http2debug=2,x509sha1=1", true},
		// net/http matches substrings. A 0 elsewhere does not turn it off.
		{"http2debug=1,http2debug=0", true},
		{"http2debug=10", true},
		// Other HTTP/2 settings leave the debug logging off.
		{"http2client=0", false},
		{"http2server=0,http2xconnect=1", false},
	}
	for _, c := range cases {
		if got := http2DebugOn(c.godebug); got != c.want {
			t.Errorf("http2DebugOn(%q) = %v, want %v", c.godebug, got, c.want)
		}
	}
}
