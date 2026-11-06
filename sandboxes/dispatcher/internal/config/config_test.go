package config

import (
	"flag"
	"strings"
	"testing"
	"time"
)

func load(t *testing.T, args ...string) (*Config, error) {
	t.Helper()
	return loadAt(t, "https://api.example/", args...)
}

// loadAt is load with baseURL as the value of ANTHROPIC_BASE_URL.
func loadAt(t *testing.T, baseURL string, args ...string) (*Config, error) {
	t.Helper()
	t.Setenv("ANTHROPIC_ENVIRONMENT_KEY", "k")
	t.Setenv("ANTHROPIC_ENVIRONMENT_ID", "env_1")
	t.Setenv("ANTHROPIC_BASE_URL", baseURL)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	cfg := Bind(fs)
	if err := fs.Parse(append([]string{"-namespace", "ns"}, args...)); err != nil {
		t.Fatal(err)
	}
	return cfg, cfg.Load()
}

func TestMalformedEnvironmentIsAnError(t *testing.T) {
	t.Setenv("CMA_DISPATCHER_API_REQUEST_TIMEOUT", "later")
	t.Setenv("CMA_DISPATCHER_HANDOFF_WAIT", "soon")
	t.Setenv("CMA_DISPATCHER_RETENTION_DAYS", "many")
	t.Setenv("CMA_DISPATCHER_SKIP_NETWORK_PROBE", "maybe")
	_, err := load(t)
	for _, name := range []string{"API_REQUEST_TIMEOUT", "HANDOFF_WAIT", "RETENTION_DAYS", "SKIP_NETWORK_PROBE"} {
		if err == nil || !strings.Contains(err.Error(), name) {
			t.Errorf("err = %v, want it to name %s", err, name)
		}
	}
}

func TestAFlagOverridesItsMalformedVariable(t *testing.T) {
	for _, tc := range []struct{ variable, flag string }{
		{"HANDOFF_WAIT", "-handoff-wait=5m"},
		{"SWEEP_INTERVAL", "-sweep-interval=1m"},
		{"RETENTION_DAYS", "-retention-days=2"},
		{"MAX_CONCURRENT_HANDOFFS", "-max-concurrent-handoffs=3"},
		{"SKIP_NETWORK_PROBE", "-skip-network-probe=true"},
	} {
		t.Run(tc.flag, func(t *testing.T) {
			t.Setenv("CMA_DISPATCHER_"+tc.variable, "nonsense")
			if _, err := load(t); err == nil || !strings.Contains(err.Error(), tc.variable) {
				t.Errorf("err = %v, want it to name %s", err, tc.variable)
			}
			if _, err := load(t, tc.flag); err != nil {
				t.Errorf("err = %v, want the flag to win over its malformed variable", err)
			}
		})
	}
	t.Setenv("CMA_DISPATCHER_HANDOFF_WAIT", "5min")
	cfg, err := load(t, "-handoff-wait", "5m")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HandoffWait != 5*time.Minute {
		t.Errorf("HandoffWait = %v, want 5m", cfg.HandoffWait)
	}
	// The flag overrides only its own variable.
	t.Setenv("CMA_DISPATCHER_RETENTION_DAYS", "many")
	_, err = load(t, "-handoff-wait", "5m")
	if err == nil || !strings.Contains(err.Error(), "RETENTION_DAYS") || strings.Contains(err.Error(), "HANDOFF_WAIT") {
		t.Errorf("err = %v, want it to name RETENTION_DAYS and not HANDOFF_WAIT", err)
	}
}

func TestEnvironmentDefaultsAndFlags(t *testing.T) {
	t.Setenv("CMA_DISPATCHER_API_REQUEST_TIMEOUT", "45s")
	t.Setenv("CMA_DISPATCHER_HANDOFF_WAIT", "90s")
	t.Setenv("CMA_DISPATCHER_MAX_CONCURRENT_HANDOFFS", "3")
	t.Setenv("CMA_DISPATCHER_SKIP_NETWORK_PROBE", "true")
	cfg, err := load(t, "-retention-days", "2")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIRequestTimeout != 45*time.Second || cfg.HandoffWait != 90*time.Second || cfg.MaxConcurrentHandoffs != 3 || !cfg.SkipNetworkProbe || cfg.Retention != 48*time.Hour {
		t.Errorf("cfg = %+v", cfg)
	}
	if _, err := load(t, "-retention-days", "-1"); err == nil {
		t.Error("a negative retention must be refused")
	}
	if _, err := load(t, "-api-request-timeout", "0"); err == nil || !strings.Contains(err.Error(), "api-request-timeout") {
		t.Errorf("err = %v, want a zero request timeout refused", err)
	}
}

func TestEmptyTokenAudienceIsAnError(t *testing.T) {
	if _, err := load(t, "-token-audience", ""); err == nil || !strings.Contains(err.Error(), "token-audience") {
		t.Errorf("err = %v, want the empty audience refused", err)
	}
}

// The workspace claim size is a Kubernetes quantity; zero turns the claim
// off, and anything else is refused at start.
func TestWorkspaceClaimSize(t *testing.T) {
	cfg, err := load(t)
	if err != nil {
		t.Fatal(err)
	}
	if got := cfg.WorkspaceClaimSize.String(); got != "10Gi" {
		t.Errorf("default -workspace-claim-size = %s, want 10Gi", got)
	}
	cfg, err = load(t, "-workspace-claim-size", "0", "-workspace-storage-class", "fast")
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.WorkspaceClaimSize.IsZero() || cfg.WorkspaceStorageClass != "fast" {
		t.Errorf("size %s class %q; want zero and fast", cfg.WorkspaceClaimSize.String(), cfg.WorkspaceStorageClass)
	}
	for _, bad := range []string{"ten gigs", "-1Gi"} {
		if _, err := load(t, "-workspace-claim-size", bad); err == nil {
			t.Errorf("-workspace-claim-size=%q must be refused", bad)
		}
	}
}

// A Secret filled by hand can hold an address no request can be sent to, with
// a user name and a password in it. The error reaches the pod's log. It names
// the variable and holds no piece of the value, in either case: Go lower-cases
// what it reads as a scheme.
func TestABaseURLWithNoHostIsRefusedAndNotShown(t *testing.T) {
	for name, value := range map[string]string{
		"a space in the password":             "http://Qz7Kx9Wv:Jq4Z p8Xy@cma.invalid:8400",
		"a bad percent escape":                "http://Qz7Kx9Wv:Jq4Z%zzp8Xy@cma.invalid:8400",
		"a control character in the password": "http://Qz7Kx9Wv:Jq4Z\tp8Xy@cma.invalid:8400",
		"a line break at the end":             "http://Qz7Kx9Wv:Jq4Zp8Xy@cma.invalid:8400\n",
		"a password where the port goes":      "http://Qz7Kx9Wv:Jq4Zp8Xy",
		"a slash in the password":             "http://Qz7Kx9Wv:Jq4Z/p8Xy@cma.invalid:8400",
		"no scheme":                           "Qz7Kx9Wv:Jq4Zp8Xy@cma.invalid:8400",
		"no scheme and a token":               "Jq4Zp8Xy@cma.invalid:8400",
		"no scheme and no credentials":        "cma.invalid:8400",
		"no host":                             "http://Qz7Kx9Wv:Jq4Zp8Xy@",
		"a scheme alone":                      "http://",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := loadAt(t, value)
			if err == nil || !strings.HasPrefix(err.Error(), "ANTHROPIC_BASE_URL is not a URL that names a host (its value is not shown)") {
				t.Fatalf("err = %v, want a refusal that names ANTHROPIC_BASE_URL", err)
			}
			for _, piece := range []string{"qz7kx9wv", "jq4z", "p8xy", "cma.invalid"} {
				if strings.Contains(strings.ToLower(err.Error()), piece) {
					t.Errorf("the error holds %q of the value: %v", piece, err)
				}
			}
		})
	}
}

// Go ends a URL's host at the first "/", "?" or "#". A password that holds one
// raw, after digits Go takes for a port, leaves a URL that parses. Its host is
// the user name, and the environment key would go there. The same rule refuses
// an "@" in a path, whose host no reader can tell from such a user name.
func TestABaseURLWithAnAtSignAfterItsHostIsRefusedAndNotShown(t *testing.T) {
	for name, value := range map[string]string{
		"a slash after digits in the password":         "http://Qz7Kx9Wv:123/Jq4Zp8Xy@cma.invalid:8400",
		"a question mark after digits in the password": "http://Qz7Kx9Wv:123?Jq4Zp8Xy@cma.invalid:8400",
		"a hash after digits in the password":          "http://Qz7Kx9Wv:123#Jq4Zp8Xy@cma.invalid:8400",
		"a slash in a token":                           "http://Qz7Kx9Wv/Jq4Zp8Xy@cma.invalid:8400",
		"a question mark in a token":                   "http://Qz7Kx9Wv?Jq4Zp8Xy@cma.invalid:8400",
		"a hash in a token":                            "http://Qz7Kx9Wv#Jq4Zp8Xy@cma.invalid:8400",
		"an at sign in the path":                       "http://cma.invalid/@Qz7Kx9Wv",
		"an at sign in the path, and a slash after it": "http://cma.invalid/@Qz7Kx9Wv/",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := loadAt(t, value)
			if err == nil || !strings.HasPrefix(err.Error(), `ANTHROPIC_BASE_URL (its value is not shown) has an "@" after a "/", "?" or "#", so its host is unclear`) {
				t.Fatalf("err = %v, want a refusal that names ANTHROPIC_BASE_URL and the at sign", err)
			}
			for _, piece := range []string{"qz7kx9wv", "123", "jq4z", "p8xy", "cma.invalid"} {
				if strings.Contains(strings.ToLower(err.Error()), piece) {
					t.Errorf("the error holds %q of the value: %v", piece, err)
				}
			}
		})
	}
}

// Every address the API client can send a request to loads, and BaseURL keeps
// its user name and password: sandboxes are handed that value.
func TestABaseURLWithAHostLoads(t *testing.T) {
	for value, want := range map[string]string{
		"https://cma.example":                            "https://cma.example",
		"https://cma.example/":                           "https://cma.example",
		"http://127.0.0.1:8400":                          "http://127.0.0.1:8400",
		"http://[::1]:8400":                              "http://[::1]:8400",
		"http://:8400":                                   "http://:8400",
		"HTTPS://cma.example/behind/a/path/":             "HTTPS://cma.example/behind/a/path",
		"http://Qz7Kx9Wv:Jq4Zp8Xy@cma.example:8400":      "http://Qz7Kx9Wv:Jq4Zp8Xy@cma.example:8400",
		"http://Qz7Kx9Wv:Jq4Z%20p8Xy@cma.example:8400":   "http://Qz7Kx9Wv:Jq4Z%20p8Xy@cma.example:8400",
		"http://Qz7Kx9Wv:Jq4Z%25zzp8Xy@cma.example:8400": "http://Qz7Kx9Wv:Jq4Z%25zzp8Xy@cma.example:8400",
		"http://Qz7Kx9Wv:1%2F2%3F3%234@cma.example:8400": "http://Qz7Kx9Wv:1%2F2%3F3%234@cma.example:8400",
		"http://Qz7Kx9Wv:Jq4Z@p8Xy@cma.example:8400":     "http://Qz7Kx9Wv:Jq4Z@p8Xy@cma.example:8400",
		"https://cma.example/%40tenant":                  "https://cma.example/%40tenant",
	} {
		cfg, err := loadAt(t, value)
		if err != nil {
			t.Errorf("%s: err = %v, want it to load", value, err)
			continue
		}
		if cfg.BaseURL != want {
			t.Errorf("%s: BaseURL = %q, want %q", value, cfg.BaseURL, want)
		}
	}
}
