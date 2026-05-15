package ratelimiter

import (
	"time"
)

type Config struct {
	Every time.Duration
	Rate  int
	Burst int
}

type Limiter struct {
	conf      *Config
	tokens    float64
	lastCheck time.Time
}

func NewLimiter(conf *Config) *Limiter {
	return &Limiter{
		conf:      conf,
		tokens:    float64(conf.Burst),
		lastCheck: time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	rate := float64(l.conf.Rate) / l.conf.Every.Seconds()

	elapsed := time.Since(l.lastCheck).Seconds()
	l.tokens += elapsed * rate
	if l.tokens > float64(l.conf.Burst) {
		l.tokens = float64(l.conf.Burst)
	}
	l.lastCheck = time.Now()

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}
