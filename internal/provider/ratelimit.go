package provider

import (
	"net/http"
	"sync"
	"time"
)

// rateLimitedTransport is an http.RoundTripper that paces outgoing requests
// to at most one every minInterval. It serializes requests so retries are
// rate-limited along with initial attempts.
type rateLimitedTransport struct {
	base        http.RoundTripper
	mu          sync.Mutex
	lastRequest time.Time
	minInterval time.Duration
}

func newRateLimitedTransport(base http.RoundTripper, rps float64) *rateLimitedTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &rateLimitedTransport{
		base:        base,
		minInterval: time.Duration(float64(time.Second) / rps),
	}
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	wait := t.minInterval - time.Since(t.lastRequest)
	if wait > 0 {
		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			t.mu.Unlock()
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}
	t.lastRequest = time.Now()
	t.mu.Unlock()
	return t.base.RoundTrip(req)
}
