package middlewares

import (
	"sync"
	"time"
)


type rateLimter struct{
	mu sync.Mutex
	visitors map[string]int
	resetTime  time.Duration
}


func