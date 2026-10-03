package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// hashLength is how much of the digest lands in the lockfile. 128 bits is far
// past what change detection needs and keeps the file readable; the actual
// integrity anchor for URL-sourced content is the pinned commit SHA, not this.
const hashLength = 32

// digest fingerprints data for the lockfile. The length prefix is part of
// every recorded hash, so dropping it would make every resource look changed.
func digest(data []byte) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d:", len(data))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))[:hashLength]
}

// canonicalJSON renders a decoded JSON value so that two bodies meaning the
// same thing hash the same. encoding/json already sorts map keys; what this
// adds is one spelling per number, since YAML's int and JSON's float64 would
// otherwise disagree.
func canonicalJSON(v any) ([]byte, error) {
	return json.Marshal(canonicalize(v))
}

func canonicalize(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[k] = canonicalize(val)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, val := range t {
			out[i] = canonicalize(val)
		}
		return out
	case json.Number:
		// Straight from a UseNumber decode. Re-render it through the same
		// rules as everything else so a server-side 4.0 and a local 4 agree.
		if i, err := t.Int64(); err == nil {
			return json.Number(strconv.FormatInt(i, 10))
		}
		if f, err := t.Float64(); err == nil {
			candidate := canonicalize(f).(json.Number)
			if candidate == t {
				return t
			}
			if f == 0 {
				// Avoid constructing enormous powers of ten for an underflowed
				// value. Only an all-zero significand is equivalent to zero.
				mantissa := strings.FieldsFunc(t.String(), func(r rune) bool { return r == 'e' || r == 'E' })[0]
				if strings.Trim(mantissa, "-0.") == "" {
					return candidate
				}
				return t
			}
			original, originalOK := new(big.Rat).SetString(t.String())
			rendered, renderedOK := new(big.Rat).SetString(candidate.String())
			// Keep float normalization only when it preserves the exact decimal
			// value. Rounding here can hide a real edit from change detection.
			if originalOK && renderedOK && original.Cmp(rendered) == 0 {
				return candidate
			}
		}
		return t
	case int, int64, uint64:
		return json.Number(fmt.Sprint(t))
	case float64:
		// A YAML `4` decodes to int, a JSON `4` to float64. Render whole
		// floats as integers so the two agree. Everything else goes through
		// the shortest round-trippable form: fixed-point would round 1e-9 to
		// zero and overflow on anything past the int64 range.
		if t == math.Trunc(t) && math.Abs(t) < 1<<53 {
			return json.Number(strconv.FormatInt(int64(t), 10))
		}
		return json.Number(strconv.FormatFloat(t, 'g', -1, 64))
	default:
		return v
	}
}

func hashBody(body map[string]any) (string, error) {
	data, err := canonicalJSON(body)
	if err != nil {
		return "", err
	}
	return digest(data), nil
}

// normalizeRemote strips server-owned fields and empty values so the remote
// fingerprint is stable across reads.
func normalizeRemote(spec KindSpec, remote map[string]any) map[string]any {
	out := deepCopyMap(remote)
	for _, path := range spec.computed {
		deletePath(out, path)
	}
	for _, path := range spec.writeOnly {
		deleteWildcardPath(out, path)
	}
	return dropEmpty(out)
}

// dropEmpty drops nulls, empty strings, and empty maps and lists, recursing
// into nested maps (including maps inside lists). The API renders "unset" as
// null in some places and as an empty array in others; treating them alike
// keeps a fingerprint from flapping.
func dropEmpty(m map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		switch t := v.(type) {
		case nil:
			continue
		case map[string]any:
			nested := dropEmpty(t)
			if len(nested) == 0 {
				continue
			}
			out[k] = nested
		case []any:
			if len(t) == 0 {
				continue
			}
			items := make([]any, 0, len(t))
			for _, item := range t {
				if nested, ok := item.(map[string]any); ok {
					items = append(items, dropEmpty(nested))
					continue
				}
				items = append(items, item)
			}
			out[k] = items
		case string:
			if t == "" {
				continue
			}
			out[k] = t
		default:
			out[k] = v
		}
	}
	return out
}
