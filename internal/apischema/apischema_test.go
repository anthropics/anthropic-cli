package apischema_test

import (
	"sort"
	"testing"

	"github.com/anthropics/anthropic-cli/v2/internal/apischema"
)

// resourceNames returns the keys of apischema.Resources in a stable order so
// failures are reported the same way on every run.
func resourceNames() []string {
	names := make([]string, 0, len(apischema.Resources))
	for name := range apischema.Resources {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Field.Path says "No two fields of a resource share a Path", and Resource.Fields
// says the fields are "sorted by Path".
func TestFieldsAreSortedByPathWithoutDuplicates(t *testing.T) {
	t.Parallel()
	for _, name := range resourceNames() {
		fields := apischema.Resources[name].Fields
		for i := 1; i < len(fields); i++ {
			prev, cur := fields[i-1].Path, fields[i].Path
			if prev == cur {
				t.Errorf("%s: path %q appears twice", name, cur)
			} else if prev > cur {
				t.Errorf("%s: path %q comes after %q, but fields are sorted by Path", name, cur, prev)
			}
		}
	}
}

// Reference says Type and IDField are set when the value is an object, and
// unset when the value is the ID itself, so one without the other is a mistake.
func TestReferencesSetTypeAndIDFieldTogether(t *testing.T) {
	t.Parallel()
	for _, name := range resourceNames() {
		for _, field := range apischema.Resources[name].Fields {
			for _, ref := range field.References {
				if ref.Resource == "" {
					t.Errorf("%s.%s: a reference names no resource", name, field.Path)
				}
				if (ref.Type == "") != (ref.IDField == "") {
					t.Errorf("%s.%s: reference to %q sets Type=%q and IDField=%q, want both or neither",
						name, field.Path, ref.Resource, ref.Type, ref.IDField)
				}
			}
		}
	}
}
