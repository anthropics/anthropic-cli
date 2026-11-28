package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatVersionRendersWhateverTypeTheDecoderChose(t *testing.T) {
	assert.Equal(t, "", FormatVersion(nil))
	assert.Equal(t, "1700000000", FormatVersion("1700000000"))
	assert.Equal(t, "7", FormatVersion(json.Number("7")))
	assert.Equal(t, "7", FormatVersion(float64(7)))
	assert.Equal(t, "1700000000000", FormatVersion(float64(1700000000000)))
	assert.Equal(t, "2.5", FormatVersion(2.5))
	assert.Equal(t, "7", FormatVersion(int64(7)))
	assert.Equal(t, "7", FormatVersion(7))
}

func TestRemoteVersionReadsTheKindsVersionField(t *testing.T) {
	spec := KindSpec{VersionField: "version"}

	assert.Equal(t, "4", remoteVersion(spec, map[string]any{"version": float64(4), "other": "x"}))
	assert.Equal(t, "", remoteVersion(spec, map[string]any{"other": "x"}))
	assert.Equal(t, "", remoteVersion(spec, nil))

	// A kind without a version field has nothing to report, even if the
	// object happens to carry a "version" key.
	assert.Equal(t, "", remoteVersion(KindSpec{}, map[string]any{"version": "3"}))
}

func TestVersionValueFollowsTheKindsVersionType(t *testing.T) {
	asInt := KindSpec{VersionIsInt: true}
	asString := KindSpec{}

	assert.Nil(t, versionValue(asInt, ""))
	assert.Nil(t, versionValue(asString, ""))

	assert.Equal(t, int64(12), versionValue(asInt, "12"))
	assert.Equal(t, "12", versionValue(asString, "12"))

	// A recorded string that is not a number goes back out unchanged.
	assert.Equal(t, "v12", versionValue(asInt, "v12"))
}

func TestCurrentVersionPrefersTheServersView(t *testing.T) {
	spec := KindSpec{VersionField: "version"}
	entry := &LockEntry{Version: "3"}

	both := &Change{spec: spec, Entry: entry, Remote: map[string]any{"version": float64(5)}}
	assert.Equal(t, "5", currentVersion(both))

	// Without a remote version the lockfile's record is used.
	noRemote := &Change{spec: spec, Entry: entry, Remote: map[string]any{}}
	assert.Equal(t, "3", currentVersion(noRemote))

	// A create has neither.
	assert.Equal(t, "", currentVersion(&Change{spec: spec}))
}
