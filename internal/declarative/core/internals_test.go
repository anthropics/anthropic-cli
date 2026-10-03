package core

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanonicalizeKeepsFloatsDistinguishable(t *testing.T) {
	// Fixed-point formatting rounded anything below 1e-6 to the same string,
	// so two different values hashed identically and a real edit looked like
	// no change.
	tiny, err := hashBody(map[string]any{"x": 1e-9})
	require.NoError(t, err)
	tinier, err := hashBody(map[string]any{"x": 1e-12})
	require.NoError(t, err)
	assert.NotEqual(t, tiny, tinier)

	huge, err := hashBody(map[string]any{"x": 1e21})
	require.NoError(t, err)
	huger, err := hashBody(map[string]any{"x": 2e21})
	require.NoError(t, err)
	assert.NotEqual(t, huge, huger)

	// Whole floats and ints still agree, which is what lets a YAML `4` and a
	// JSON `4.0` compare equal.
	a, _ := hashBody(map[string]any{"x": 4})
	b, _ := hashBody(map[string]any{"x": 4.0})
	assert.Equal(t, a, b)
}

func TestEquivalentToleratesServerNormalization(t *testing.T) {
	tests := []struct {
		name string
		want any
		got  any
		same bool
	}{
		{"identical strings", "a", "a", true},
		{"a string is not an object without a declared shorthand", "x",
			map[string]any{"id": "x"}, false},
		{"server fills in extra keys",
			map[string]any{"type": "limited"},
			map[string]any{"type": "limited", "allow_mcp_servers": false}, true},
		{"declared key differs",
			map[string]any{"type": "limited"},
			map[string]any{"type": "unrestricted"}, false},
		{"null and empty list are both unset", nil, []any{}, true},
		{"list length matters", []any{"a"}, []any{"a", "b"}, false},
		{"numeric forms agree", int64(4), float64(4), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.same, equivalent(tc.want, tc.got))
		})
	}
}

// A field declared with a Shorthand compares its scalar form equal to the
// object the server answers with; without that the plan would report a change
// on every run and never converge.
func TestShorthandFieldComparesEqualToItsObjectForm(t *testing.T) {
	spec := NewRegistry(testSchema{}, KindSpec{
		Kind:   "k",
		Fields: Fields{"model": {Shorthand: "id"}},
		Build:  func(*Candidate) (map[string]any, error) { return nil, nil },
	}).specOrZero("k")

	same := &Change{spec: spec,
		Desired: map[string]any{"model": "m-5", "name": "n"},
		Remote:  map[string]any{"model": map[string]any{"id": "m-5", "speed": "standard"}, "name": "n"},
	}
	assert.Nil(t, same.buildDiff())

	moved := &Change{spec: spec,
		Desired: map[string]any{"model": "m-6"},
		Remote:  map[string]any{"model": map[string]any{"id": "m-5"}},
	}
	d := moved.buildDiff()
	require.NotNil(t, d)
	require.Contains(t, d.Fields, "model")
	assert.Equal(t, DiffText, d.Fields["model"].Kind)
	assert.Equal(t, "m-5", d.Fields["model"].Before, "compared against the member the scalar stands for")
	assert.Equal(t, "m-6", d.Fields["model"].After, "the plan shows what the file says, not the expansion")
}

// Inserting one element into a list must read as one addition. Pairing by
// position instead would report every later element as changed and bury the
// actual edit.
func TestListElementsPairByIdentityBeforePosition(t *testing.T) {
	spec := NewRegistry(testSchema{}, KindSpec{
		Kind:   "k",
		Fields: Fields{"parts": {MatchBy: []string{"sku"}}},
		Build:  func(*Candidate) (map[string]any, error) { return nil, nil },
	}).specOrZero("k")

	c := &Change{spec: spec,
		Remote: map[string]any{
			"parts": []any{
				map[string]any{"sku": "b", "qty": 1},
				map[string]any{"sku": "c", "qty": 1},
			},
			"tags": []any{map[string]any{"name": "x"}, map[string]any{"name": "y"}},
		},
		Desired: map[string]any{
			"parts": []any{
				map[string]any{"sku": "a", "qty": 1}, // inserted at the front
				map[string]any{"sku": "b", "qty": 1}, // moved, unchanged
				map[string]any{"sku": "c", "qty": 2}, // moved and edited
			},
			"tags": []any{map[string]any{"name": "y"}}, // x dropped; default keys pair y by name
		},
	}
	d := c.buildDiff()
	require.NotNil(t, d)

	parts := d.Fields["parts"]
	require.Equal(t, DiffList, parts.Kind)
	require.Len(t, parts.Items, 3)
	assert.Equal(t, ItemDiff{Before: -1, After: 0, Diff: &Diff{Kind: DiffAdded, After: map[string]any{"sku": "a", "qty": 1}}}, parts.Items[0])
	assert.Equal(t, ItemDiff{Before: 0, After: 1}, parts.Items[1], "b only moved")
	assert.Equal(t, 1, parts.Items[2].Before)
	assert.Equal(t, 2, parts.Items[2].After)
	require.NotNil(t, parts.Items[2].Diff)
	assert.Contains(t, parts.Items[2].Diff.Fields, "qty", "c is diffed field by field, not replaced")

	tags := d.Fields["tags"]
	require.Len(t, tags.Items, 2)
	assert.Equal(t, DiffRemoved, tags.Items[0].Diff.Kind)
	assert.Equal(t, -1, tags.Items[0].After)
	assert.Equal(t, ItemDiff{Before: 1, After: 0}, tags.Items[1], "y pairs by name and only moved")
}

// A secret nested inside a list element must not surface anywhere in the tree,
// whichever way the element is reported — added, removed, or changed.
func TestSecretsInsideListElementsNeverReachTheDiff(t *testing.T) {
	spec := NewRegistry(testSchema{}, KindSpec{
		Kind:   "k",
		Fields: Fields{"resources[].token": {WriteOnly: true}, "resources": {MatchBy: []string{"url"}}},
		Build:  func(*Candidate) (map[string]any, error) { return nil, nil },
	}).specOrZero("k")
	const secret = "ghp_SUPERSECRET"

	c := &Change{spec: spec,
		Sensitive: map[string]bool{"resources[].token": true},
		Remote: map[string]any{"resources": []any{
			map[string]any{"url": "a", "token": secret + "-old"}, // dropped
			map[string]any{"url": "b"},                           // token never returned
		}},
		Desired: map[string]any{"resources": []any{
			map[string]any{"url": "b", "token": secret},          // kept
			map[string]any{"url": "c", "token": secret + "-new"}, // added
		}},
	}
	d := c.buildDiff()
	require.NotNil(t, d)
	dump, err := json.Marshal(d)
	require.NoError(t, err)
	assert.NotContains(t, string(dump), secret)

	items := d.Fields["resources"].Items
	require.Len(t, items, 3)
	assert.Equal(t, DiffSensitive, items[0].Diff.Kind, "a removed element holding a secret is withheld whole")
	assert.Equal(t, DiffWriteOnly, items[1].Diff.Fields["token"].Kind, "b's token cannot be compared, only re-sent")
	assert.Equal(t, DiffSensitive, items[2].Diff.Kind, "an added element holding a secret is withheld whole")
}

func TestPayloadResponsePrefersTheVersionJustMinted(t *testing.T) {
	// A payload-backed resource has two version fields in play: the object's
	// own (here "latest_version") and the one an upload response reports for
	// the version it just created. A response holding both must not resolve to
	// the older of the two.
	// A payload-backed kind is one with no Build: its content is not a body.
	bundled := NewRegistry(testSchema{},
		KindSpec{Kind: "bundle", VersionField: "latest_version"}).specOrZero("bundle")
	assert.Equal(t, "1002", responseVersion(bundled,
		map[string]any{"latest_version": "1001", "version": "1002"}))
	assert.Equal(t, "1001", responseVersion(bundled, map[string]any{"latest_version": "1001"}))
	assert.Equal(t, "1002", responseVersion(bundled, map[string]any{"version": "1002"}))

	// An ordinary kind reads its own field and nothing else.
	plain := NewRegistry(testSchema{},
		KindSpec{Kind: "plain", VersionField: "version", Build: func(*Candidate) (map[string]any, error) { return nil, nil }}).specOrZero("plain")
	assert.Equal(t, "4", responseVersion(plain, map[string]any{"version": int64(4)}))
	none := NewRegistry(testSchema{},
		KindSpec{Kind: "none", Build: func(*Candidate) (map[string]any, error) { return nil, nil }}).specOrZero("none")
	assert.Empty(t, responseVersion(none, map[string]any{"version": "ignored"}))
}

func TestSplitFrontmatter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantFront string
		wantBody  string
		wantErr   string
	}{
		{
			name:      "frontmatter and body",
			input:     "---\nmodel: sonnet\n---\nHello\n",
			wantFront: "model: sonnet\n",
			wantBody:  "Hello\n",
		},
		{
			name:     "no frontmatter is all body",
			input:    "Just a prompt\n",
			wantBody: "Just a prompt\n",
		},
		{
			name:      "empty frontmatter",
			input:     "---\n---\nbody",
			wantFront: "",
			wantBody:  "body",
		},
		{
			name:      "body may contain a horizontal rule",
			input:     "---\nname: x\n---\nintro\n\n---\n\nmore\n",
			wantFront: "name: x\n",
			wantBody:  "intro\n\n---\n\nmore\n",
		},
		{
			name:      "leading BOM is tolerated",
			input:     "\ufeff---\nname: x\n---\nbody",
			wantFront: "name: x\n",
			wantBody:  "body",
		},
		{
			name:    "unterminated frontmatter is an error",
			input:   "---\nname: x\nbody without a closing fence\n",
			wantErr: "never closed",
		},
		{
			name:     "a triple dash mid-line is not a fence",
			input:    "text --- more\n",
			wantBody: "text --- more\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			front, body, err := splitFrontmatter([]byte(tc.input))
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantFront, string(front))
			assert.Equal(t, tc.wantBody, string(body))
		})
	}
}

// TestYAMLKeysAreAlwaysStrings pins the decoder behaviour ParseYAMLMap relies
// on: keys YAML would happily type as bool, int, null or a timestamp still come
// back as strings, at every depth, so a body is always json.Marshal-able. If a
// dependency bump breaks this, it breaks here rather than as an opaque marshal
// error inside hashBody.
func TestYAMLKeysAreAlwaysStrings(t *testing.T) {
	const doc = `
on: yes
1: one
true: yes it is
null: nothing
2020-01-01: a date
nested:
  2: two
listed:
  - 3: three
tagged: !!map {4: four}
`
	got, err := ParseYAMLMap([]byte(doc), "probe")
	require.NoError(t, err)

	canonical, err := canonicalJSON(got)
	require.NoError(t, err, "a parsed body must survive the marshal in canonicalJSON")
	assert.JSONEq(t, `{
		"on": "yes", "1": "one", "true": "yes it is", "null": "nothing",
		"2020-01-01": "a date", "nested": {"2": "two"},
		"listed": [{"3": "three"}], "tagged": {"4": "four"}
	}`, string(canonical))
}

func TestExactJSONNumbersRemainDistinctInFingerprints(t *testing.T) {
	pairs := [][2]string{
		{"18446744073709551614", "18446744073709551615"},
		{"-9223372036854775809", "-9223372036854775810"},
		{"1.000000000000000001", "1.000000000000000002"},
		{"1e-400", "2e-400"},
		{"1e-999999999", "0"},
	}
	for _, pair := range pairs {
		t.Run(pair[0], func(t *testing.T) {
			bodies := make([]map[string]any, 2)
			for i, number := range pair {
				decoder := json.NewDecoder(strings.NewReader(`{"nested":{"values":[` + number + `]}}`))
				decoder.UseNumber()
				require.NoError(t, decoder.Decode(&bodies[i]))
			}
			before, err := json.Marshal(bodies)
			require.NoError(t, err)
			first, err := hashBody(bodies[0])
			require.NoError(t, err)
			second, err := hashBody(bodies[1])
			require.NoError(t, err)
			assert.NotEqual(t, first, second, "a real numeric edit must change its fingerprint")
			after, err := json.Marshal(bodies)
			require.NoError(t, err)
			assert.Equal(t, before, after, "canonicalization must not mutate its inputs")
		})
	}
}

func TestExactNumberNormalizationKeepsRepresentableAliases(t *testing.T) {
	groups := [][]any{
		{4, int64(4), float64(4), json.Number("4"), json.Number("4.0"), json.Number("4e0")},
		{0, float64(0), json.Number("-0.0"), json.Number("0e-999999999")},
		{0.1, json.Number("0.1000"), json.Number("1e-1")},
		{float64(1e20), json.Number("100000000000000000000"), json.Number("1e20")},
	}
	for _, group := range groups {
		expected, err := hashBody(map[string]any{"value": group[0]})
		require.NoError(t, err)
		for _, value := range group[1:] {
			actual, err := hashBody(map[string]any{"value": value})
			require.NoError(t, err)
			assert.Equal(t, expected, actual, "equivalent numeric forms: %v", value)
		}
	}
}

func TestExactRemoteNumberEditIsDetectedAsDrift(t *testing.T) {
	registry := testRegistry()
	root := writeTree(t, map[string]string{"widgets/value.json": `{"name":"value","amount":1}`})
	path := filepath.Join(root, "widgets/value.json")
	loader := NewLoader(registry, root, nil)
	require.NoError(t, loader.Add(context.Background(), []string{path}))
	source := loader.Sources()["./widgets/value.json"]
	require.NotNil(t, source)
	decode := func(value string) map[string]any {
		decoder := json.NewDecoder(strings.NewReader(`{"name":"value","amount":` + value + `}`))
		decoder.UseNumber()
		var remote map[string]any
		require.NoError(t, decoder.Decode(&remote))
		return remote
	}
	oldRemote := decode("1.000000000000000001")
	oldHash, err := hashBody(source.Body)
	require.NoError(t, err)
	remoteHash, err := hashBody(oldRemote)
	require.NoError(t, err)
	lock := newLockfile(filepath.Join(root, testLockfileName))
	lock.Resources[source.Key] = &LockEntry{Kind: "widget", ID: "wdgt_1", Hash: oldHash, RemoteHash: remoteHash}
	require.NoError(t, lock.Save())
	restored, err := LoadLockfile(registry, lock.Path)
	require.NoError(t, err)
	planner := Planner{Registry: registry, Lock: restored}
	builder := bodyBuilder{registry: registry, slots: loader.slots}
	before, err := planner.planOne(builder, source, map[string]Target{}, remoteResult{Object: oldRemote})
	require.NoError(t, err)
	assert.Equal(t, ActionNoop, before.Action)
	updatedRemote := decode("1.000000000000000002")
	change, err := planner.planOne(builder, source, map[string]Target{}, remoteResult{Object: updatedRemote})
	require.NoError(t, err)
	assert.Equal(t, ActionUpdate, change.Action)
	assert.True(t, change.Drift)
	assert.ErrorContains(t, change.Blocked, "changed outside this config")
	require.NotNil(t, change.Diff)
	require.Contains(t, change.Diff.Fields, "amount")
	assert.Equal(t, json.Number("1.000000000000000002"), change.Remote["amount"])
}
