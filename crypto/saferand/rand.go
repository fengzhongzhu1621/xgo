package saferand

import (
	"math/rand"
	"sync"
)

// SafeRand is the safe random functions struct.
type SafeRand struct {
	r  *rand.Rand
	mu sync.Mutex
}

// NewSafeRand creates a SafeRand using the given seed.
func NewSafeRand(seed int64) *SafeRand {
	c := &SafeRand{
		r: rand.New(rand.NewSource(seed)),
	}
	return c
}

// Int63n provides a random int64.
func (c *SafeRand) Int63n(n int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := c.r.Int63n(n)
	return res
}

// Intn provides a random int.
func (c *SafeRand) Intn(n int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := c.r.Intn(n)
	return res
}

// Float64 provides a random float64.
func (c *SafeRand) Float64() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := c.r.Float64()
	return res
}
