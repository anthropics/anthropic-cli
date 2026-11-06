// Package preflight checks an install before the dispatcher runs: the Agent
// Sandbox API, the sandbox runtime, storage, namespace policy, the
// dispatcher Deployment, the environment key, and the dispatcher's RBAC.
package preflight

import (
	"context"
	"errors"
	"fmt"
	"io"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// Options names the objects the checks look at.
type Options struct {
	Namespace                string
	Shape                    sandbox.Shape
	SandboxTemplateName      string
	WarmPoolName             string
	DispatcherServiceAccount string
	DeploymentName           string
	TokenAudience            string
	// ServiceHost is the host:port sandbox pods reach the dispatcher at. The
	// network check names it in the steps the operator runs by hand.
	ServiceHost      string
	SkipNetworkProbe bool
	// WorkspaceClaim is the claim the dispatcher gives each direct Sandbox's
	// /workspace; nil means an emptyDir a suspend discards.
	WorkspaceClaim *sandbox.WorkspaceClaim
}

// EnvironmentVerifier reports whether the environment key reads the queue.
type EnvironmentVerifier interface {
	VerifyEnvironment(ctx context.Context) error
}

// Checker runs the checks with the operator's own cluster credentials.
type Checker struct {
	Kube    kubernetes.Interface
	Dyn     dynamic.Interface
	API     EnvironmentVerifier
	Options Options
}

type check struct {
	name string
	run  func(ctx context.Context) error
}

// ErrWarning marks a finding that does not fail the run.
var ErrWarning = errors.New("warning")

// Run prints one line per check to out and returns the number of failures.
func (c *Checker) Run(ctx context.Context, out io.Writer) int {
	failures := 0
	for _, ck := range c.checks() {
		err := ck.run(ctx)
		switch {
		case err == nil:
			_, _ = fmt.Fprintf(out, "ok    %s\n", ck.name)
		case errors.Is(err, ErrWarning):
			_, _ = fmt.Fprintf(out, "warn  %s: %v\n", ck.name, err)
		default:
			failures++
			_, _ = fmt.Fprintf(out, "FAIL  %s: %v\n", ck.name, err)
		}
	}
	return failures
}

func (c *Checker) checks() []check {
	if c.Options.Shape.Direct() {
		return []check{
			{"agent sandbox API is served", c.apiServed},
			{"sandbox runtime is gVisor or Kata", c.runtimeClass},
			{"sandbox template turns the managed network policy off", c.managedPolicyOff},
			{"sandbox template carries no kubelet credential and limits every container", c.directContainers},
			{"workspace survives a suspend", c.directWorkspace},
			{"the add-on admits a sandbox built from the template", c.admitsSandbox},
			{"network policy isolates sandboxes", c.networkProbe},
			{"namespace enforces restricted pod security", c.podSecurity},
			{"dispatcher runs as one replica with Recreate", c.deployment},
			{"environment key reads the work queue", c.environmentKey},
			{"dispatcher service account holds its role", c.rbac},
		}
	}
	return []check{
		{"agent sandbox API is served", c.apiServed},
		{"warm pool makes its pods from the checked template", c.warmPoolTemplate},
		{"sandbox runtime is gVisor or Kata", c.runtimeClass},
		{"sandbox template turns the managed network policy off", c.managedPolicyOff},
		{"sandbox template keeps the token in the fetch container", c.containers},
		{"network policy isolates sandboxes", c.networkProbe},
		{"workspace storage class exists", c.storageClass},
		{"namespace enforces restricted pod security", c.podSecurity},
		{"dispatcher runs as one replica with Recreate", c.deployment},
		{"environment key reads the work queue", c.environmentKey},
		{"dispatcher service account holds its role", c.rbac},
	}
}
