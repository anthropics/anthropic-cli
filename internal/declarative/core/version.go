package core

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Versions arrive from the server in whatever type the decoder chose, are
// recorded in the lockfile as strings, and go back out in the form each kind's
// API expects.

// remoteVersion extracts the comparable version from a fetched object, as the
// string the lockfile records whatever numeric type the decoder chose for it.
func remoteVersion(spec KindSpec, remote map[string]any) string {
	field := spec.VersionField
	if remote == nil || field == "" {
		return ""
	}
	return FormatVersion(remote[field])
}

// FormatVersion renders a version as the string the lockfile records, whatever
// type the decoder chose for it.
func FormatVersion(version any) string {
	switch v := version.(type) {
	case nil:
		return ""
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

// responseVersion reads the version out of a write response.
//
// A payload upload answers with the version it just created, under `version`;
// the resource object itself reports its current one under the kind's own
// version field. Where a response holds both, the freshly minted one is the
// truthful answer, so it wins rather than being a fallback.
func responseVersion(spec KindSpec, obj map[string]any) string {
	if spec.payloadBacked {
		if v, ok := obj["version"].(string); ok && v != "" {
			return v
		}
	}
	return remoteVersion(spec, obj)
}

// versionValue renders a recorded version string in the form the API expects
// for that kind: some kinds take versions as integers, others as opaque strings.
func versionValue(spec KindSpec, version string) any {
	if version == "" {
		return nil
	}
	if spec.VersionIsInt {
		if n, err := strconv.ParseInt(version, 10, 64); err == nil {
			return n
		}
	}
	return version
}

// currentVersion prefers the server's view over the lockfile's.
func currentVersion(change *Change) string {
	if v := remoteVersion(change.spec, change.Remote); v != "" {
		return v
	}
	if change.Entry != nil {
		return change.Entry.Version
	}
	return ""
}
