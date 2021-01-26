package util

import (
	"sync"
	"time"
)

type IPRateLimiter struct {
	ips          map[string]*Bucket
	mu           *sync.RWMutex
	fillInterval time.Duration
	quantum      int64
	capacity     int64
}

func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		ips:          make(map[string]*Bucket),
		mu:           &sync.RWMutex{},
		fillInterval: 60 * time.Second,
		quantum:      60,
		capacity:     60,
	}
}

func (ipRL *IPRateLimiter) AddIP(ip string) *Bucket {
	ipRL.mu.Lock()
	defer ipRL.mu.Unlock()

	limiter := NewBucket(ipRL.fillInterval, ipRL.quantum, ipRL.capacity)
	ipRL.ips[ip] = limiter

	return limiter
}

func (ipRL *IPRateLimiter) GetLimiter(ip string) *Bucket {
	ipRL.mu.Lock()
	limiter, exists := ipRL.ips[ip]

	if !exists {
		ipRL.mu.Unlock()
		return ipRL.AddIP(ip)
	}

	ipRL.mu.Unlock()

	return limiter
}
