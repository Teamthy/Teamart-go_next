package websocket

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket per connection.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[int64]*tokenBucket
	rate    int
	burst   int
}

type tokenBucket struct {
	tokens        int
	lastTokenTime time.Time
}

// NewRateLimiter creates a new websocket rate limiter.
func NewRateLimiter(rate int, burst int) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[int64]*tokenBucket),
		rate:    rate,
		burst:   burst,
	}
}

// Allow returns whether the given user is allowed to send a message.
func (rl *RateLimiter) Allow(userID int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok := rl.buckets[userID]
	if !ok {
		bucket = &tokenBucket{tokens: rl.burst, lastTokenTime: time.Now()}
		rl.buckets[userID] = bucket
	}

	elapsed := time.Since(bucket.lastTokenTime)
	replenish := int(elapsed.Seconds()) * rl.rate
	if replenish > 0 {
		bucket.tokens = min(bucket.tokens+replenish, rl.burst)
		bucket.lastTokenTime = time.Now()
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
