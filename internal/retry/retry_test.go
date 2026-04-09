package retry

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMiddleware_NoRetryOnSuccess(t *testing.T) {
	callCount := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 3
	middleware := NewMiddleware(config)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(r)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestNewMiddleware_RetryOn429(t *testing.T) {
	callCount := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 5
	config.BaseDelay = 10 * time.Millisecond
	middleware := NewMiddleware(config)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(r)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestNewMiddleware_RetryOn500(t *testing.T) {
	callCount := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 2
	config.BaseDelay = 10 * time.Millisecond
	middleware := NewMiddleware(config)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(r)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
	if callCount != 3 { // initial + 2 retries
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestNewMiddleware_NoRetryOn400(t *testing.T) {
	callCount := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusBadRequest)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 3
	middleware := NewMiddleware(config)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(r)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestNewMiddleware_RetryAfterHeader(t *testing.T) {
	callCount := 0
	start := time.Now()
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 3
	middleware := NewMiddleware(config)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(r)
	})

	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
	// Should have waited at least 1 second due to Retry-After
	if elapsed < 900*time.Millisecond {
		t.Errorf("expected delay of ~1s, got %v", elapsed)
	}
}

func TestNewMiddleware_MaxRetriesExceeded(t *testing.T) {
	callCount := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	config := DefaultConfig()
	config.MaxRetries = 2
	config.BaseDelay = 10 * time.Millisecond
	middleware := NewMiddleware(config)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return http.DefaultClient.Do(r)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", resp.StatusCode)
	}
	if callCount != 3 { // initial + 2 retries
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestCalculateDelay_ExponentialBackoff(t *testing.T) {
	config := DefaultConfig()
	config.BaseDelay = 100 * time.Millisecond
	config.MaxDelay = 1 * time.Second

	delays := make([]time.Duration, 5)
	for i := 1; i <= 5; i++ {
		delays[i-1] = calculateDelay(i, config, nil)
	}

	// Check that delays increase exponentially
	for i := 1; i < len(delays); i++ {
		if delays[i] < delays[i-1] {
			t.Errorf("delay %d (%v) should be >= delay %d (%v)", i+1, delays[i], i, delays[i-1])
		}
	}

	// Check max delay cap
	if calculateDelay(10, config, nil) > config.MaxDelay+100*time.Millisecond {
		t.Error("delay should be capped at MaxDelay + jitter")
	}
}

func TestCalculateDelay_RetryAfterSeconds(t *testing.T) {
	config := DefaultConfig()
	resp := &http.Response{
		Header: http.Header{"Retry-After": []string{"5"}},
	}

	delay := calculateDelay(1, config, resp)
	if delay != 5*time.Second {
		t.Errorf("expected 5s delay, got %v", delay)
	}
}

func TestCalculateDelay_RetryAfterDate(t *testing.T) {
	config := DefaultConfig()
	future := time.Now().Add(2 * time.Second)
	// Use http.TimeFormat for proper HTTP-date format (RFC 7231 IMF-fixdate)
	resp := &http.Response{
		Header: http.Header{"Retry-After": []string{future.UTC().Format(http.TimeFormat)}},
	}

	delay := calculateDelay(1, config, resp)
	// Should be approximately 2 seconds (allow wider range for test stability)
	if delay < 1400*time.Millisecond || delay > 2600*time.Millisecond {
		t.Errorf("expected ~2s delay, got %v", delay)
	}
}
