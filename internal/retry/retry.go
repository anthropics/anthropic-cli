package retry

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Middleware is a function that wraps an HTTP request
type Middleware = func(*http.Request, func(*http.Request) (*http.Response, error)) (*http.Response, error)

// Config holds retry configuration
type Config struct {
	MaxRetries  int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	ShouldRetry func(resp *http.Response, err error) bool
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   30 * time.Second,
		ShouldRetry: func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			// Retry on rate limit or server errors
			return resp.StatusCode == http.StatusTooManyRequests ||
				(resp.StatusCode >= 500 && resp.StatusCode < 600)
		},
	}
}

// NewMiddleware creates a retry middleware with the given config
func NewMiddleware(config Config) Middleware {
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		var lastResp *http.Response
		var lastErr error

		for attempt := 0; attempt <= config.MaxRetries; attempt++ {
			if attempt > 0 {
				delay := calculateDelay(attempt, config, lastResp)
				select {
				case <-time.After(delay):
					// Continue to retry
				case <-req.Context().Done():
					return nil, req.Context().Err()
				}
			}

			resp, err := next(req)
			lastResp = resp
			lastErr = err

			if !config.ShouldRetry(resp, err) {
				return resp, err
			}

			// Close response body to prevent leaks
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}

		// All retries exhausted
		if lastErr != nil {
			return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
		}
		return lastResp, nil
	}
}

// calculateDelay determines how long to wait before retrying
func calculateDelay(attempt int, config Config, lastResp *http.Response) time.Duration {
	// Check for Retry-After header
	if lastResp != nil {
		if retryAfter := lastResp.Header.Get("Retry-After"); retryAfter != "" {
			// Try parsing as seconds first
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				return time.Duration(seconds) * time.Second
			}
			// Try parsing as HTTP-date
			if date, err := http.ParseTime(retryAfter); err == nil {
				delay := time.Until(date)
				if delay > 0 {
					return delay
				}
			}
		}
	}

	// Exponential backoff with jitter: min(2^attempt * baseDelay, maxDelay) + jitter
	exponential := config.BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
	if exponential > config.MaxDelay {
		exponential = config.MaxDelay
	}

	// Add jitter (0-100ms)
	jitter := time.Duration(rand.Intn(100)) * time.Millisecond
	return exponential + jitter
}
