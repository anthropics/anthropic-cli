package preflight

import (
	"strings"
	"testing"
)

// The probe is not implemented, so running it fails the run. Only
// -skip-network-probe turns the line into a warning. Either line names the
// dispatcher host in the steps the operator runs by hand.
func TestPreflightNetworkCheck(t *testing.T) {
	for name, tc := range map[string]struct {
		probe bool
		want  string
	}{
		"skipped": {false, skippedProbe},
		"run":     {true, "FAIL  network policy isolates sandboxes: not implemented in this version; " + byHand},
	} {
		c := healthy()
		c.networkProbe = tc.probe
		wantLine(t, name, run(t, c), tc.want)
	}
}

// A Deployment whose pod template names no serviceAccountName runs as the
// namespace's default ServiceAccount. The deployment check compares that
// account, not the empty name, with -dispatcher-service-account, so an
// install run as default passes on both paths.
func TestPreflightDeploymentWithoutServiceAccountRunsAsDefault(t *testing.T) {
	c := healthy()
	c.deploymentSA, c.dispatcherSA = "", "default"
	if out := run(t, c); strings.Contains(out, "FAIL") {
		t.Errorf("upstream: a deployment naming no service account, checked as default, fails:\n%s", out)
	}
	d := healthyDirect()
	d.deploymentSA, d.dispatcherSA = "", "default"
	if out := runDirect(t, d); strings.Contains(out, "FAIL") {
		t.Errorf("add-on: a deployment naming no service account, checked as default, fails:\n%s", out)
	}
}
