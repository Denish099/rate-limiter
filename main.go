package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	rate       float64
	capacity   float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	elapsed := now.Sub(tb.lastRefill).Seconds()

	tb.tokens += elapsed * tb.rate

	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	tb.lastRefill = now

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}
func main() {
	limiter := NewTokenBucket(1.0, 3.0)

	fmt.Println("Sending burst of 4 requests...")
	for i := 1; i <= 4; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED\n", i)
		} else {
			fmt.Printf("Request %d: DROPPED (Rate Limited)\n", i)
		}
	}

	fmt.Println("\nWaiting 2 seconds...")
	time.Sleep(2 * time.Second)

	fmt.Println("Sending another request...")
	if limiter.Allow() {
		fmt.Println("Request 5: ALLOWED")
	} else {
		fmt.Println("Request 5: DROPPED")
	}
}
