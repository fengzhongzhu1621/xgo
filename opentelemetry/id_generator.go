package opentelemetry

import (
	"math/rand"
	"sync"
)

// randomIDGenerator generates random span ID.
type randomIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

// newSpanID returns a non-negative span ID randomly.
func (gen *randomIDGenerator) newSpanID() SpanID {
	gen.Lock()
	defer gen.Unlock()

	return SpanID(gen.randSource.Int63())
}

// newRandomIDGenerator returns a random ID generator with the given seed.
// if seed equals zero, the seed will be set to the current time.
func newRandomIDGenerator(seed int64) *randomIDGenerator {
	return &randomIDGenerator{randSource: rand.New(rand.NewSource(seed))}
}
