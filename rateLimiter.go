package main

import (
	"golang.org/x/time/rate"
	"sync"
)

type IPRateLimiter struct {
	ips         map[string]*rate.Limiter
	mu          *sync.RWMutex
	RateLimit   rate.Limit
	TokenBucket int
}

func NewIPRateLimiter(rateLimit rate.Limit, tokenBucket int) *IPRateLimiter {
	ipRateLimiter := &IPRateLimiter{
		ips:         make(map[string]*rate.Limiter),
		mu:          &sync.RWMutex{},
		RateLimit:   rateLimit,
		TokenBucket: tokenBucket,
	}

	return ipRateLimiter
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	//We need to lock the Struct because various requests may want to write/delete/read here at the same time.
	i.mu.Lock()
	// It's a Good practice to defer the unlock. There is a minimum penalty around 10 microseconds.
	defer i.mu.Unlock()
	// create a new limiter.
	limiter := rate.NewLimiter(i.RateLimit, i.TokenBucket)
	// Save the limiter in the hashmap for the IP associated.
	i.ips[ip] = limiter

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	//defer i.mu.Unlock()
	limiter, limiterExists := i.ips[ip]
	if !limiterExists {
		// If limiter already exists for the given IP return the existing limiter.
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	return limiter
}

func (i *IPRateLimiter) GetNextAvailableTime(limiter *rate.Limiter){
		limiter.Limit()
}