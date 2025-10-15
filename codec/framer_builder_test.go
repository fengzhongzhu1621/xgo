package codec

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

var _ ISafeFramer = (*fakeSafeFramer)(nil)

type fakeSafeFramer struct {
	safe bool
}

func (f *fakeSafeFramer) ReadFrame() ([]byte, error) {
	return nil, nil
}

func (f *fakeSafeFramer) IsSafe() bool {
	return f.safe
}

func TestIsSafeFramer(t *testing.T) {
	safeFrame := fakeSafeFramer{safe: true}
	assert.Equal(t, true, IsSafeFramer(&safeFrame))

	noSafeFrame := fakeSafeFramer{}
	assert.Equal(t, false, IsSafeFramer(&noSafeFrame))

	assert.Equal(t, false, IsSafeFramer(10))
}
