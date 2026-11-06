package preflight

import (
	"testing"
)

// The role check asks the cluster about every verb the dispatcher uses on
// each path, so an install whose role lacks any one of them fails the check
// by name. Each path's list here is the expected role, spelled out.
func TestPreflightRBACRequiresEveryVerb(t *testing.T) {
	type role struct {
		group, resource string
		verbs           []string
	}
	upstream := []role{
		{"extensions.agents.x-k8s.io", "sandboxclaims", []string{"create", "get", "list", "watch", "delete", "patch"}},
		{"extensions.agents.x-k8s.io", "sandboxwarmpools", []string{"get"}},
		{"agents.x-k8s.io", "sandboxes", []string{"get", "list", "watch", "patch"}},
		{"", "pods", []string{"get", "list", "watch"}},
	}
	direct := []role{
		{"agents.x-k8s.io", "sandboxes", []string{"create", "get", "list", "watch", "patch", "delete"}},
		{"extensions.agents.x-k8s.io", "sandboxtemplates", []string{"get"}},
		{"", "secrets", []string{"create"}},
		{"", "pods", []string{"get", "list", "watch"}},
	}
	want := func(r role, verb string) string {
		return "FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher may not " + verb + " " + r.resource + groupSuffix(r.group)
	}
	for _, r := range upstream {
		for _, verb := range r.verbs {
			c := healthy()
			c.denyAccess = access{r.resource, verb}
			wantLine(t, "upstream role lacking "+verb+" on "+r.resource, run(t, c), want(r, verb))
		}
	}
	for _, r := range direct {
		for _, verb := range r.verbs {
			c := healthyDirect()
			c.denyAccess = access{r.resource, verb}
			wantLine(t, "add-on role lacking "+verb+" on "+r.resource, runDirect(t, c), want(r, verb))
		}
	}
}
