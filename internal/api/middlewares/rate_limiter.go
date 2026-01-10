package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)


type rateLimter struct{
	mu sync.Mutex
	visitors map[string]int
	limit int
	resetTime  time.Duration
}


func NewRateLimiter(limit int , resetTime time.Duration) *rateLimter{
	rl := &rateLimter{
		visitors: make(map[string]int),
		limit: limit,
		resetTime: resetTime,
	}
	 go rl.resetVisitorCount()
	return rl
}


func (rl *rateLimter ) resetVisitorCount(){
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()

	}
}


func (rl *rateLimter) Middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       rl.mu.Lock()
	   defer rl.mu.Unlock()
	   visitorIP := r.RemoteAddr  // you might want to extract the IP in more sophiscated way
	   rl.visitors[visitorIP]++
	   fmt.Printf("Visitor counr from IP %v is %v", visitorIP , rl.visitors[visitorIP])
     if rl.visitors[visitorIP] > rl.limit {
		http.Error(w, "Too many Request", http.StatusTooManyRequests)
		return 

	}
	

		next.ServeHTTP(w, r)
	})
}