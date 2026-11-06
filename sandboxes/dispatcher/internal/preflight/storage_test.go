package preflight

import (
	"strings"
	"testing"
)

func TestPreflightDefaultStorageClass(t *testing.T) {
	for _, annotation := range []string{"storageclass.kubernetes.io/is-default-class", "storageclass.beta.kubernetes.io/is-default-class"} {
		c := healthy()
		c.storageClass = ""
		c.defaultClassAnnotation = annotation
		if out := run(t, c); strings.Contains(out, "FAIL") {
			t.Errorf("a claim template naming no class should take the class marked by %s:\n%s", annotation, out)
		}
	}
}
