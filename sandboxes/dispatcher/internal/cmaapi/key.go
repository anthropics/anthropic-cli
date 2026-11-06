package cmaapi

import (
	"errors"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

// withoutKey returns err with a placeholder in place of the environment key.
// The SDK's error quotes a refusal's body and request id, and callers log it.
// A proxy on the way to the API can copy the request it refuses into either.
func withoutKey(err error) error {
	var apiErr *anthropic.Error
	if !errors.As(err, &apiErr) {
		return err
	}
	key := strings.TrimPrefix(apiErr.Request.Header.Get("Authorization"), "Bearer ")
	if text := err.Error(); key != "" && strings.Contains(text, key) {
		return errors.New(strings.ReplaceAll(text, key, "[environment key]"))
	}
	return err
}
