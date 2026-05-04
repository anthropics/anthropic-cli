package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFederationIdentityTokenFunc covers the assertion resolver that the CLI
// hands to option.WithFederationTokenProvider: literal value, file (re-read
// each call, whitespace-trimmed, empty rejected), mutual-exclusion, and the
// nothing-configured nil-nil case.
func TestFederationIdentityTokenFunc(t *testing.T) {
	t.Run("literal assertion", func(t *testing.T) {
		fn, err := federation{Assertion: "literal-jwt"}.IdentityTokenFunc()
		require.NoError(t, err)
		require.NotNil(t, fn)
		tok, err := fn(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "literal-jwt", tok)
	})

	t.Run("file is re-read and trimmed", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "assertion.jwt")
		require.NoError(t, os.WriteFile(path, []byte("  first-jwt\n"), 0600))
		fn, err := federation{AssertionFile: path}.IdentityTokenFunc()
		require.NoError(t, err)
		require.NotNil(t, fn)
		tok, err := fn(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "first-jwt", tok)
		// Rotate; second call must pick up the new value.
		require.NoError(t, os.WriteFile(path, []byte("second-jwt"), 0600))
		tok, err = fn(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "second-jwt", tok)
	})

	t.Run("empty file rejected", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.jwt")
		require.NoError(t, os.WriteFile(path, []byte("  \n"), 0600))
		fn, err := federation{AssertionFile: path}.IdentityTokenFunc()
		require.NoError(t, err)
		_, err = fn(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is empty")
	})

	t.Run("literal and file are mutually exclusive", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "a.jwt")
		require.NoError(t, os.WriteFile(path, []byte("x"), 0600))
		_, err := federation{Assertion: "literal", AssertionFile: path}.IdentityTokenFunc()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mutually exclusive")
	})

	t.Run("nothing configured returns nil, nil", func(t *testing.T) {
		fn, err := federation{}.IdentityTokenFunc()
		require.NoError(t, err)
		assert.Nil(t, fn)
	})
}

func TestFederationMissing(t *testing.T) {
	for _, tc := range []struct {
		name string
		fed  federation
		want []string
	}{
		{"empty → nil (fall through)", federation{}, nil},
		{"complete → nil", federation{Assertion: "j", Rule: "r", OrganizationID: "o"}, nil},
		{"only rule set", federation{Rule: "r"},
			[]string{
				"--identity-token[-file] / ANTHROPIC_IDENTITY_TOKEN[_FILE]",
				"--organization-id / ANTHROPIC_ORGANIZATION_ID",
			}},
		{"only assertion set", federation{Assertion: "j"},
			[]string{
				"--federation-rule / ANTHROPIC_FEDERATION_RULE_ID",
				"--organization-id / ANTHROPIC_ORGANIZATION_ID",
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.fed.Missing())
		})
	}
}
