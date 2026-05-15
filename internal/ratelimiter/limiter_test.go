package ratelimiter

import (
	"testing"
	"time"
)

// TestLimiterAllow tests the token bucket rate limiting algorithm.
// Expected behavior: Must allow requests up to the rate, then deny until tokens replenish.
func TestLimiterAllow(t *testing.T) {
	limiter := NewLimiter(&Config{
		Every: time.Second,
		Rate:  10,
		Burst: 5,
	})

	// Should allow 5 requests immediately (burst)
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed (within burst), but was denied", i)
		}
	}

	// 6th request should be denied (no tokens left)
	if limiter.Allow() {
		t.Error("Request should be denied (no tokens left), but was allowed")
	}

	// Wait for some tokens to replenish (0.1 seconds = 1 token at 10/sec)
	time.Sleep(150 * time.Millisecond)

	// Should now allow one more request
	if !limiter.Allow() {
		t.Error("Request should be allowed after token replenishment, but was denied")
	}
}

// TestLimiterBurst tests that burst size is respected.
// Expected behavior: Initial requests can exceed rate, up to burst limit.
func TestLimiterBurst(t *testing.T) {
	conf := &Config{
		Every: time.Second,
		Rate:  10,
		Burst: 5,
	}

	limiter := NewLimiter(conf)

	// Should allow exactly 'burst' requests initially
	for i := 0; i < conf.Burst; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed (within burst of %d), but was denied", i, conf.Burst)
		}
	}

	// (burst + 1)th request should be denied
	if limiter.Allow() {
		t.Error("Request beyond burst should be denied, but was allowed")
	}
}

// TestLimiterReplenishment tests token replenishment over time.
// Expected behavior: Tokens should replenish at specified rate per second.
func TestLimiterReplenishment(t *testing.T) {
	conf := &Config{
		Every: time.Second,
		Rate:  10,
		Burst: 3,
	}

	limiter := NewLimiter(conf)

	// Use the initial token
	for range conf.Burst {
		if !limiter.Allow() {
			t.Error("Initial request should be allowed")
		}
	}

	// Next request should be denied (no tokens)
	if limiter.Allow() {
		t.Error("Request should be denied (no tokens)")
	}

	// Wait for 2 tokens to replenish (0.2 seconds at 10 req/sec)
	time.Sleep(210 * time.Millisecond)

	// Should now allow requests
	if !limiter.Allow() {
		t.Error("Request should be allowed after replenishment")
	}

	if !limiter.Allow() {
		t.Error("Second request should be allowed after replenishment")
	}

	// Third should be denied
	if limiter.Allow() {
		t.Error("Third request should be denied (not enough replenishment)")
	}
}

// TestLimiterBurstCap tests that token count respects burst maximum.
// Expected behavior: Tokens should not exceed burst limit, even after long idle periods.
func TestLimiterBurstCap(t *testing.T) {
	conf := &Config{
		Every: time.Second,
		Rate:  10,
		Burst: 3,
	}

	limiter := NewLimiter(conf)

	// Wait for a long time (should not accumulate beyond burst)
	time.Sleep(1200 * time.Millisecond)

	// Should still only have 'burst' tokens available
	for i := 0; i < conf.Burst; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should be allowed (within burst cap of %d)", i, conf.Burst)
		}
	}

	// Next request should be denied
	if limiter.Allow() {
		t.Error("Request beyond burst cap should be denied")
	}
}

// TestNewLimiter tests limiter initialization.
// Expected behavior: New limiter should start with burst tokens available.
func TestNewLimiter(t *testing.T) {
	conf := &Config{
		Every: time.Second,
		Rate:  5,
		Burst: 10,
	}

	limiter := NewLimiter(conf)

	if limiter.conf.Rate != conf.Rate {
		t.Errorf("Limiter conf.Rate = %d, want %d", limiter.conf.Rate, conf.Rate)
	}

	if limiter.conf.Burst != conf.Burst {
		t.Errorf("Limiter conf.Burst = %d, want %d", limiter.conf.Burst, conf.Burst)
	}

	if limiter.tokens != float64(conf.Burst) {
		t.Errorf("Limiter initial tokens = %.1f, want %d", limiter.tokens, conf.Burst)
	}
}
