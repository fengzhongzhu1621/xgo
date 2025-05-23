package ratelimiter

import "github.com/juju/ratelimit"

// RateLimiter TODO
type RateLimiter interface {
	// TryAccept returns true if a token is taken immediately. Otherwise,
	// it returns false.
	TryAccept() bool

	// Accept will wait and not return unless a token becomes available.
	Accept()

	// QPS returns QPS of this rate limiter
	QPS() int64

	// Burst returns the burst of this rate limiter
	Burst() int64

	// AcceptMany will wait and not return unless the token becomes available.
	AcceptMany(count int64)
}

// NewRateLimiter TODO
func NewRateLimiter(qps, burst int64) RateLimiter {
	limiter := ratelimit.NewBucketWithRate(float64(qps), burst)
	return &tokenBucket{
		limiter: limiter,
		qps:     qps,
		burst:   burst,
	}
}

type tokenBucket struct {
	limiter *ratelimit.Bucket
	qps     int64
	burst   int64
}

// TryAccept TODO
func (t *tokenBucket) TryAccept() bool {
	return t.limiter.TakeAvailable(1) == 1
}

// Accept TODO
func (t *tokenBucket) Accept() {
	t.limiter.Wait(1)
}

// QPS TODO
func (t *tokenBucket) QPS() int64 {
	return t.qps
}

// Burst TODO
func (t *tokenBucket) Burst() int64 {
	return t.burst
}

// AcceptMany accept many token
func (t *tokenBucket) AcceptMany(count int64) {
	t.limiter.Wait(count)
}

// NewMockRateLimiter TODO
func NewMockRateLimiter() RateLimiter {
	return &mockRatelimiter{}
}

type mockRatelimiter struct{}

// TryAccept TODO
func (*mockRatelimiter) TryAccept() bool {
	return true
}

// Accept TODO
func (*mockRatelimiter) Accept() {

}

// QPS TODO
func (*mockRatelimiter) QPS() int64 {
	return 0
}

// Burst TODO
func (*mockRatelimiter) Burst() int64 {
	return 0
}

// AcceptMany accept many token
func (*mockRatelimiter) AcceptMany(count int64) {
}
