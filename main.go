package main

import (
    "fmt"
    "sync"
    "time"
)

type RateLimiter struct {
    rate       int           
    capacity   int
    tokens     int
    mu         sync.Mutex
    ticker     *time.Ticker
    stopTicker chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, capacity int) *RateLimiter {
    rl := &RateLimiter{
        rate:       rate,
        capacity:   capacity,
        tokens:     capacity, 
        ticker:     time.NewTicker(time.Second), 
        stopTicker: make(chan struct{}),
    }

    go rl.startReplenishment()
    return rl
}

// startReplenishment adds tokens at a fixed rate
func (rl *RateLimiter) startReplenishment() {
    for {
        select {
        case <-rl.ticker.C:
            rl.mu.Lock()
            if rl.tokens < rl.capacity {
                rl.tokens++
            }
            rl.mu.Unlock()
        case <-rl.stopTicker:
            return
        }
    }
}

// Allow checks if a request can be made
func (rl *RateLimiter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

		fmt.Println("Current tokens: ", rl.tokens)

    if rl.tokens > 0 {
        rl.tokens--
        return true
    }

    return false
}

// Stop stops the token replenishment
func (rl *RateLimiter) Stop() {
    rl.ticker.Stop()
    close(rl.stopTicker)
}

func main() {
    rateLimiter := NewRateLimiter(5, 5)

    defer rateLimiter.Stop()

    for i := 0; i < 15; i++ {
        if rateLimiter.Allow() {
            fmt.Println("Request", i+1, "allowed")
        } else {
            fmt.Println("Request", i+1, "denied")
        }

        time.Sleep(200 * time.Millisecond) 
    }
}
