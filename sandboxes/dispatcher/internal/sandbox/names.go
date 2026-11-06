package sandbox

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

const (
	// LabelSessionID carries the CMA session id on the session's SandboxClaim.
	LabelSessionID = "anthropic.com/session-id"
	// LabelEnvironmentID carries the CMA environment id on the claim, so
	// dispatchers of different environments sharing a namespace keep apart.
	LabelEnvironmentID = "anthropic.com/environment-id"
	// LabelManagedBy marks the claims this dispatcher owns.
	LabelManagedBy = "app.kubernetes.io/managed-by"
	// ManagedByValue is the LabelManagedBy value the dispatcher writes.
	ManagedByValue = "cma-dispatcher"
	// AnnotationLastActive holds the RFC 3339 time of the claim's last handoff.
	AnnotationLastActive = "anthropic.com/last-active"
	// AnnotationHandoffRestarts holds the worker's restart count at the
	// claim's last handoff, so a dispatcher started later judges the
	// worker's exit by count rather than by clock.
	AnnotationHandoffRestarts = "anthropic.com/handoff-restarts"
	// AnnotationHandoffPodUID holds the pod the claim's last handoff went
	// to, so a dispatcher started later never mistakes a replacement pod
	// for the handed one.
	AnnotationHandoffPodUID = "anthropic.com/handoff-pod-uid"

	maxDNSLabel = 63
	hashHexLen  = 12
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// ClaimName derives the SandboxClaim name for sessionID (the upstream path).
func ClaimName(sessionID string) string { return objectName(sessionID) }

// SandboxName derives the Sandbox name for sessionID (the add-on path, where
// the dispatcher creates the Sandbox directly instead of through a claim).
func SandboxName(sessionID string) string { return objectName(sessionID) }

// NonceSecretName names the Secret holding sandboxName's assignment nonce.
// The sandbox's fetch container reads it through an env valueFrom, so the
// nonce appears in no Sandbox or Pod spec.
func NonceSecretName(sandboxName string) string { return sandboxName + "-nonce" }

// RearmedNonceSecretName names a fresh Secret for re-arming sandboxName's
// nonce. The dispatcher holds no Secrets verb but create, so a re-arm cannot
// reuse a name; the random suffix makes each one new.
func RearmedNonceSecretName(sandboxName string) (string, error) {
	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		return "", fmt.Errorf("naming a fresh nonce secret for sandbox %s: %w", sandboxName, err)
	}
	return NonceSecretName(sandboxName) + "-" + hex.EncodeToString(suffix), nil
}

// objectName derives the DNS-1123 label naming a session's sandbox object:
// the lowercased id with each run of non-alphanumerics collapsed to "-", cut
// so the name plus twelve hex digits of sha256(sessionID) fits the label.
func objectName(sessionID string) string {
	sum := sha256.Sum256([]byte(sessionID))
	suffix := hex.EncodeToString(sum[:hashHexLen/2])
	base := strings.Trim(nonAlphanumeric.ReplaceAllString(strings.ToLower(sessionID), "-"), "-")
	if base == "" {
		base = "session"
	}
	if limit := maxDNSLabel - 1 - len(suffix); len(base) > limit {
		base = strings.TrimRight(base[:limit], "-")
	}
	return base + "-" + suffix
}
