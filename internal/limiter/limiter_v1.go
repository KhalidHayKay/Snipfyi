package limiter

import "time"

type TokenBucket struct {
	Tokens     float64
	Capacity   float64
	Rate       float64
	lastRefill time.Time
}

var buckets = make(map[string]*TokenBucket)

func refill(bucket *TokenBucket) {
	timeElasped := time.Since(bucket.lastRefill).Seconds()

	bucket.Tokens += timeElasped * bucket.Rate

	if bucket.Tokens > bucket.Capacity {
		bucket.Tokens = bucket.Capacity
	}

	bucket.lastRefill = time.Now()
}

func Allow(apiKey string) bool {
	bucket, exists := buckets[apiKey]

	if !exists {
		bucket = &TokenBucket{
			Tokens:     10,
			Capacity:   10,
			Rate:       10 / 60.0,
			lastRefill: time.Now(),
		}

		buckets[apiKey] = bucket
	}

	refill(bucket)

	if bucket.Tokens > 0 {
		bucket.Tokens--
		return true // Process goes through here
	}

	return false // Process goes through here
}
