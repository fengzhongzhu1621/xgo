package backoff

import (
	"errors"
	"time"
)

// CustomizedBackoff wraps an user defined bf function.
type CustomizedBackoff struct {
	bf func(attempt int) time.Duration
}

// NewCustomizedBackoff create a new CustomizedBackoff.
func NewCustomizedBackoff(bf func(attempt int) time.Duration) (*CustomizedBackoff, error) {
	if bf == nil {
		return nil, errors.New("provided bf function must not be nil")
	}

	return &CustomizedBackoff{bf: bf}, nil
}

// backoff simply calls the wrapped bf function.
func (bf *CustomizedBackoff) Backoff(attempt int) time.Duration {
	return bf.bf(attempt)
}
