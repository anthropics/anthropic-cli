package sandbox

import (
	"regexp"
	"strings"
	"testing"
)

var dnsLabel = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

func TestClaimName(t *testing.T) {
	long := "session_" + strings.Repeat("X", 100)
	cases := []struct{ in, want string }{
		{"session_01ABC", "session-01abc-" + ClaimName("session_01ABC")[len("session-01abc-"):]},
		{"", "session-" + ClaimName("")[len("session-"):]},
		{"--Weird__ID!!", "weird-id-" + ClaimName("--Weird__ID!!")[len("weird-id-"):]},
	}
	for _, c := range cases {
		got := ClaimName(c.in)
		if got != c.want {
			t.Errorf("ClaimName(%q) = %q, want %q", c.in, got, c.want)
		}
		if !dnsLabel.MatchString(got) || len(got) > 63 {
			t.Errorf("ClaimName(%q) = %q is not a DNS-1123 label", c.in, got)
		}
	}
	if a, b := ClaimName(long), ClaimName(long+"Y"); len(a) != 63 || !dnsLabel.MatchString(a) || a == b {
		t.Errorf("long ids: %q vs %q — want 63-char distinct labels", a, b)
	}
	if ClaimName("session_01abc") == ClaimName("session_01ABC") {
		t.Error("ids differing only in case must not share a name")
	}
}
