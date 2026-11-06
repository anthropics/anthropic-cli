package preflight

import (
	"strings"
	"testing"
)

func TestPreflightHealthy(t *testing.T) {
	for _, handler := range []string{"runsc", "gvisor", "kata-qemu"} {
		c := healthy()
		c.handler = handler
		wantLine(t, "handler "+handler, run(t, c), skippedProbe)
	}
}

func TestPreflightDirectHealthy(t *testing.T) {
	if out := runDirect(t, healthyDirect()); strings.Contains(out, "FAIL") {
		t.Errorf("a healthy add-on install fails:\n%s", out)
	}
}

// The dispatcher exits on the count Run returns, so it must equal the number
// of FAIL lines on both paths; a warning counts for nothing.
func TestPreflightRunCountsFailures(t *testing.T) {
	if out, n := runWithCount(t, healthy()); n != 0 {
		t.Errorf("healthy cluster: Run returned %d, want 0:\n%s", n, out)
	}
	if out, n := runDirectWithCount(t, healthyDirect()); n != 0 {
		t.Errorf("healthy add-on install: Run returned %d, want 0:\n%s", n, out)
	}
	c := healthy()
	c.handler, c.psaLabel = "runc", "baseline"
	if out, n := runWithCount(t, c); n != 2 || strings.Count(out, "FAIL  ") != 2 {
		t.Errorf("two failing checks: Run returned %d, want 2:\n%s", n, out)
	}
	d := healthyDirect()
	d.handler, d.psaLabel = "runc", "baseline"
	if out, n := runDirectWithCount(t, d); n != 2 || strings.Count(out, "FAIL  ") != 2 {
		t.Errorf("two failing add-on checks: Run returned %d, want 2:\n%s", n, out)
	}
}
