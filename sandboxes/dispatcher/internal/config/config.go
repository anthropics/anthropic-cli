// Package config reads the dispatcher's settings: the three ANTHROPIC_*
// variables the environment Secret provides, and CMA_DISPATCHER_* variables
// (or the matching flags, which win) for everything Kubernetes-side.
package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	envPrefix     = "CMA_DISPATCHER_"
)

type Config struct {
	// EnvironmentKey authorizes every CMA API call; the environment Secret provides it.
	EnvironmentKey string
	// EnvironmentID is the environment whose work queue the dispatcher serves.
	EnvironmentID string
	// BaseURL is the CMA API the dispatcher and its sandboxes talk to.
	BaseURL string

	// Namespace holds the sandboxes; inside a pod it defaults to the pod's own.
	K8SNamespace string
	// APIVersion is the agents.x-k8s.io version the cluster serves: v1alpha1 (GKE add-on) or v1beta1 (upstream).
	K8SAPIVersion string
	// SandboxTemplate names the SandboxTemplate session pods are made from.
	K8SSandboxTemplateName string
	// WarmPool names the SandboxWarmPool new claims check out from.
	K8SWarmPoolName string
	// SandboxServiceAccount is the account sandbox pods run as; only its tokens are accepted.
	K8SSandboxServiceAccount string
	// DispatcherServiceAccount is the account whose RBAC preflight checks.
	K8SDispatcherServiceAccount string
	// DeploymentName is the dispatcher Deployment preflight checks.
	K8SDeploymentName string
	// TokenAudience is the audience a sandbox pod's projected token must carry.
	TokenAudience string
	// Listen is the address the HTTPS server binds.
	Listen string
	// ServiceHost is the host:port sandbox pods reach the dispatcher at; preflight's network check names it.
	ServiceHost string
	// TLSCertFile and TLSKeyFile are the HTTPS server's certificate and key.
	TLSCertFile string
	TLSKeyFile  string
	// APIRequestTimeout is how long one CMA API request may take before it fails.
	APIRequestTimeout time.Duration
	// HandoffWait is how long a session has to reach a running pod that collects its work item; the queue's reclaim window is twice it.
	HandoffWait time.Duration
	// SweepInterval is how often the sweep checks sandboxes against the API.
	SweepInterval time.Duration
	// Retention deletes suspended sandboxes idle for longer than it; zero keeps them.
	Retention time.Duration
	// MaxConcurrentHandoffs bounds the sessions provisioned at once; zero sizes it from the warm pool.
	MaxConcurrentHandoffs int
	// WorkspaceClaimSize is the per-session /workspace claim on the add-on path; zero backs /workspace with an emptyDir instead.
	WorkspaceClaimSize resource.Quantity
	// WorkspaceStorageClass names the StorageClass of that claim; empty takes the cluster default.
	WorkspaceStorageClass string
	// LogFormat is text or json.
	LogFormat string
	// SkipNetworkProbe turns preflight's network policy check into a warning.
	SkipNetworkProbe bool

	retentionDays      int
	workspaceClaimSize string
	// fs is the FlagSet Bind registered on. validate reads from it which
	// flags the command line set.
	fs *flag.FlagSet
	// envErrors holds each malformed CMA_DISPATCHER_* value by the flag that
	// overrides it.
	envErrors []envError
}

type envError struct {
	flag string
	err  error
}

// Bind registers the flags on fs with CMA_DISPATCHER_* values as defaults.
func Bind(fs *flag.FlagSet) *Config {
	c := &Config{fs: fs}
	fs.StringVar(&c.K8SNamespace, "namespace", env("NAMESPACE", ""), "namespace of the sandboxes (default: the service-account mount)")
	fs.StringVar(&c.K8SAPIVersion, "api-version", env("API_VERSION", "v1beta1"), "agents.x-k8s.io API version: v1alpha1 (GKE add-on) or v1beta1 (upstream)")
	fs.StringVar(&c.K8SSandboxTemplateName, "sandbox-template", env("SANDBOX_TEMPLATE", "cma-sandbox"), "SandboxTemplate name")
	fs.StringVar(&c.K8SWarmPoolName, "warm-pool", env("WARM_POOL", "cma-sandbox"), "SandboxWarmPool name")
	fs.StringVar(&c.K8SSandboxServiceAccount, "sandbox-service-account", env("SANDBOX_SERVICE_ACCOUNT", "cma-sandbox"), "ServiceAccount the sandbox pods run as")
	fs.StringVar(&c.K8SDispatcherServiceAccount, "dispatcher-service-account", env("DISPATCHER_SERVICE_ACCOUNT", "cma-dispatcher"), "ServiceAccount the dispatcher runs as (preflight)")
	fs.StringVar(&c.K8SDeploymentName, "deployment", env("DEPLOYMENT", "cma-dispatcher"), "dispatcher Deployment name (preflight)")
	fs.StringVar(&c.TokenAudience, "token-audience", env("TOKEN_AUDIENCE", "cma-dispatcher"), "audience of the projected tokens sandbox pods present")
	fs.StringVar(&c.Listen, "listen", env("LISTEN", ":8443"), "address the HTTPS server listens on")
	fs.StringVar(&c.ServiceHost, "service-host", env("SERVICE_HOST", "cma-dispatcher:8443"), "host:port sandbox pods reach the dispatcher at (preflight)")
	fs.StringVar(&c.TLSCertFile, "tls-cert-file", env("TLS_CERT_FILE", "/etc/cma-dispatcher/tls/tls.crt"), "TLS certificate file")
	fs.StringVar(&c.TLSKeyFile, "tls-key-file", env("TLS_KEY_FILE", "/etc/cma-dispatcher/tls/tls.key"), "TLS private key file")
	fs.DurationVar(&c.APIRequestTimeout, "api-request-timeout", envValue(c, "api-request-timeout", 30*time.Second, time.ParseDuration, "a duration"), "time one CMA API request may take before it fails")
	fs.DurationVar(&c.HandoffWait, "handoff-wait", envValue(c, "handoff-wait", 120*time.Second, time.ParseDuration, "a duration"), "time a session gets to reach a running pod that collects its work item; the queue's reclaim window is twice this")
	fs.DurationVar(&c.SweepInterval, "sweep-interval", envValue(c, "sweep-interval", 5*time.Minute, time.ParseDuration, "a duration"), "how often sandboxes are checked against the API")
	fs.IntVar(&c.retentionDays, "retention-days", envValue(c, "retention-days", 0, strconv.Atoi, "an integer"), "delete suspended sandboxes idle for this many days (0 keeps them)")
	fs.StringVar(&c.workspaceClaimSize, "workspace-claim-size", env("WORKSPACE_CLAIM_SIZE", "10Gi"), "add-on path: size of each session's /workspace claim (0 backs /workspace with an emptyDir a suspend discards)")
	fs.StringVar(&c.WorkspaceStorageClass, "workspace-storage-class", env("WORKSPACE_STORAGE_CLASS", ""), "add-on path: StorageClass of that claim (default: the cluster's default class)")
	fs.IntVar(&c.MaxConcurrentHandoffs, "max-concurrent-handoffs", envValue(c, "max-concurrent-handoffs", 0, strconv.Atoi, "an integer"), "sessions provisioned at once (0: the warm pool's replicas, or 8 when it cannot be read)")
	fs.StringVar(&c.LogFormat, "log-format", env("LOG_FORMAT", "text"), "log format: text or json")
	fs.BoolVar(&c.SkipNetworkProbe, "skip-network-probe", envValue(c, "skip-network-probe", false, strconv.ParseBool, "a boolean"), "preflight: skip the network policy probe")
	return c
}

// Load fills the API settings from the environment and resolves defaults.
func (c *Config) Load() error {
	c.EnvironmentKey = os.Getenv("ANTHROPIC_ENVIRONMENT_KEY")
	c.EnvironmentID = os.Getenv("ANTHROPIC_ENVIRONMENT_ID")
	c.BaseURL = strings.TrimRight(os.Getenv("ANTHROPIC_BASE_URL"), "/")
	c.Retention = time.Duration(c.retentionDays) * 24 * time.Hour
	if c.K8SNamespace == "" {
		raw, err := os.ReadFile(namespaceFile)
		if err != nil {
			return errors.New("-namespace is required outside a pod (no service-account namespace file)")
		}
		c.K8SNamespace = strings.TrimSpace(string(raw))
	}
	return c.validate()
}

func (c *Config) validate() error {
	if err := c.malformedEnv(); err != nil {
		return err
	}
	var missing []string
	for name, value := range map[string]string{
		"ANTHROPIC_ENVIRONMENT_KEY": c.EnvironmentKey,
		"ANTHROPIC_ENVIRONMENT_ID":  c.EnvironmentID,
		"ANTHROPIC_BASE_URL":        c.BaseURL,
	} {
		if value == "" {
			missing = append(missing, name)
		}
	}
	slices.Sort(missing)
	if len(missing) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	// The error holds neither the value nor url.Parse's error, which quotes
	// it. The value can hold a password, and the error reaches the pod's log.
	if u, err := url.Parse(c.BaseURL); err != nil || u.Host == "" {
		return errors.New("ANTHROPIC_BASE_URL is not a URL that names a host (its value is not shown); set it to the server's address, for example https://api.anthropic.com, with no space or line break in it and any symbol in a user name or password percent-encoded")
	}
	// Go reads the start of a user name or password as the host of such a
	// value. Requests would carry the environment key to that host.
	if atAfterHost(c.BaseURL) {
		return errors.New(`ANTHROPIC_BASE_URL (its value is not shown) has an "@" after a "/", "?" or "#", so its host is unclear; in a user name or password, write "/" as %2F, "?" as %3F and "#" as %23; in the path, query or fragment, write "@" as %40`)
	}
	if c.LogFormat != "text" && c.LogFormat != "json" {
		return fmt.Errorf("-log-format must be text or json, not %q", c.LogFormat)
	}
	if c.APIRequestTimeout <= 0 || c.HandoffWait <= 0 || c.SweepInterval <= 0 {
		return errors.New("-api-request-timeout, -handoff-wait and -sweep-interval must be positive")
	}
	if c.Retention < 0 || c.MaxConcurrentHandoffs < 0 {
		return errors.New("-retention-days and -max-concurrent-handoffs must not be negative")
	}
	if c.TokenAudience == "" {
		return errors.New("-token-audience must not be empty")
	}
	size, err := resource.ParseQuantity(c.workspaceClaimSize)
	if err != nil || size.Sign() < 0 {
		return fmt.Errorf("-workspace-claim-size must be a Kubernetes quantity such as 10Gi, or 0, not %q", c.workspaceClaimSize)
	}
	c.WorkspaceClaimSize = size
	return nil
}

// atAfterHost reports whether raw holds an "@" after the first "/", "?" or "#"
// that follows its "//". Go's url.Parse ends the host at that byte.
func atAfterHost(raw string) bool {
	_, rest, _ := strings.Cut(raw, "//")
	i := strings.IndexAny(rest, "/?#")
	return i >= 0 && strings.Contains(rest[i:], "@")
}

// malformedEnv joins the malformed variables whose flags the command line
// did not set. A flag that was set wins over its variable. The variable's
// value is then never used, and its parse failure is not an error.
func (c *Config) malformedEnv() error {
	if len(c.envErrors) == 0 {
		return nil
	}
	set := map[string]bool{}
	c.fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	var errs []error
	for _, e := range c.envErrors {
		if !set[e.flag] {
			errs = append(errs, e.err)
		}
	}
	return errors.Join(errs...)
}

func env(name, fallback string) string {
	if v := os.Getenv(envPrefix + name); v != "" {
		return v
	}
	return fallback
}

// envValue keeps a malformed variable as an error for Load instead of quietly
// using the fallback. The variable is flagName upper-cased with underscores,
// so it matches only a flag named the way Bind names them.
func envValue[T any](c *Config, flagName string, fallback T, parse func(string) (T, error), kind string) T {
	name := envPrefix + strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
	v := os.Getenv(name)
	if v == "" {
		return fallback
	}
	parsed, err := parse(v)
	if err != nil {
		c.envErrors = append(c.envErrors, envError{flag: flagName, err: fmt.Errorf("%s: %q is not %s", name, v, kind)})
		return fallback
	}
	return parsed
}
