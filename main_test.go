package main

import (
    "testing"
    "time"
)

// TestRateLimiter_Allow tests that a rate limiter correctly allows and denies requests.
// It first tests that requests are allowed when tokens are available, and then tests
// that a request is denied after all tokens are exhausted.
func TestRateLimiter_Allow(t *testing.T) {
    rl := NewRateLimiter(2, 5) // 2 tokens per second, capacity of 5
    defer rl.Stop()

    // Test allowing requests when tokens are available
    for i := 0; i < 5; i++ {
        if !rl.Allow() {
            t.Errorf("Expected request %d to be allowed, but it was denied", i+1)
        }
    }

    // Now all tokens should be exhausted
    if rl.Allow() {
        t.Errorf("Expected request to be denied, but it was allowed")
    }
}

// TestRateLimiter_ExhaustTokens tests the case where all tokens are exhausted, and the next
// request should be denied.
func TestRateLimiter_ExhaustTokens(t *testing.T) {
    rl := NewRateLimiter(1, 3) // 1 token per second, capacity of 3
    defer rl.Stop()

    // Use up all tokens
    for i := 0; i < 3; i++ {
        if !rl.Allow() {
            t.Errorf("Expected request %d to be allowed, but it was denied", i+1)
        }
    }

    // No tokens left, next request should be denied
    if rl.Allow() {
        t.Errorf("Expected request to be denied due to exhausted tokens, but it was allowed")
    }
}

// TestRateLimiter_TokenReplenishment tests that the rate limiter replenishes
// tokens over time. It verifies that after exhausting all tokens, the
// rate limiter waits for the replenishment interval to elapse before
// replenishing its tokens and allowing the next request.
func TestRateLimiter_TokenReplenishment(t *testing.T) {
    rl := NewRateLimiter(1, 2) // 1 token per second, capacity of 2
    if rl == nil {
        t.Errorf("NewRateLimiter returned a nil pointer")
    }

    defer rl.Stop()

    // Use up all tokens
    for i := 0; i < 2; i++ {
        if !rl.Allow() {
            t.Errorf("Expected request %d to be allowed, but it was denied", i+1)
        }
    }

    if rl.Allow() {
        t.Errorf("Expected request to be denied due to exhausted tokens, but it was allowed")
    }

    // Wait for replenishment (1 second for 1 token)
    time.Sleep(1 * time.Second)

    // Now the next request should be allowed
    if !rl.Allow() {
        t.Errorf("Expected request to be allowed after replenishment, but it was denied")
    }
}

// TestRateLimiter_BurstCapacity tests that the rate limiter allows requests
// up to its burst capacity. It verifies that the rate limiter can handle
// a burst of requests equal to its capacity and denies subsequent requests
// once the tokens are exhausted.
func TestRateLimiter_BurstCapacity(t *testing.T) {
    rl := NewRateLimiter(2, 3) // 2 tokens per second, capacity of 3
    defer rl.Stop()

    // We should be able to make 3 requests immediately
    for i := 0; i < 3; i++ {
        if !rl.Allow() {
            t.Errorf("Expected request %d to be allowed, but it was denied", i+1)
        }
    }

    // All tokens are exhausted now
    if rl.Allow() {
        t.Errorf("Expected request to be denied due to exhausted tokens, but it was allowed")
    }
}

func TestRateLimiter_Stop(t *testing.T) {
    rl := NewRateLimiter(1, 2) // 1 token per second, capacity of 2

    // Use up all tokens
    rl.Allow()
    rl.Allow()

    // Stop replenishment
    rl.Stop()

    // Wait to ensure that no tokens are replenished after stop
    time.Sleep(2 * time.Second)

    // Should still deny requests since replenishment stopped
    if rl.Allow() {
        t.Errorf("Expected request to be denied after stopping replenishment, but it was allowed")
    }
}
