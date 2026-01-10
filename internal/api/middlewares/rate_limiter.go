package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// rateLimter manages request quotas per visitor.
// Future Reference: Fixed-Window Algorithm
// This implementation resets the entire map at a set interval. While efficient,
// it can allow a "burst" of traffic at the very end of one window and the 
// start of the next.
type rateLimter struct {
	mu         sync.Mutex     // Logic Block: Prevents race conditions during map access
	visitors   map[string]int // Logic Block: Stores request counts keyed by IP
	limit      int            // Logic Block: Max allowed requests per window
	resetTime  time.Duration  // Logic Block: Length of the time window
}

// NewRateLimiter initializes the struct and kicks off the background cleanup.
func NewRateLimiter(limit int, resetTime time.Duration) *rateLimter {
	rl := &rateLimter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	// Logic Block: Background Goroutine
	// Future Reference: Concurrency & Background Tasks
	// Starting this in a goroutine ensures the map is cleared periodically 
	// without blocking the main server execution.
	go rl.resetVisitorCount()
	return rl
}

// resetVisitorCount clears all visitor data after every resetTime interval.
func (rl *rateLimter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		// Logic Block: Mutex Locking
		// We lock the mutex to ensure no request is trying to increment 
		// the map while we are overwriting it.
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

// Middleware checks the visitor's quota before allowing the request to proceed.
func (rl *rateLimter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		// Logic Block: Defer Unlock
		// Future Reference: Mutex Defer
		// Using defer ensures the mutex is unlocked even if the code panics 
		// or returns early, preventing "deadlocks."
		defer rl.mu.Unlock()

		// Logic Block: Identity Identification
		// r.RemoteAddr usually contains "IP:Port". In production, you might 
		// use the "X-Forwarded-For" header if behind a proxy like Nginx.
		visitorIP := r.RemoteAddr 
		rl.visitors[visitorIP]++

		fmt.Printf("Visitor count from IP %v is %v\n", visitorIP, rl.visitors[visitorIP])

		// Logic Block: Quota Enforcement
		if rl.visitors[visitorIP] > rl.limit {
			// Future Reference: HTTP 429
			// StatusTooManyRequests (429) is the standard RFC response for rate limiting.
			http.Error(w, "Too many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}