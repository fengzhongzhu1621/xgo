package throttle

// Noop defines a trivial throttle which does nothing.
type Noop struct{}

// NewNoop create a new Noop throttle.
func NewNoop() *Noop {
	return &Noop{}
}

// Allow Always return true
func (*Noop) Allow() bool { return true }

// OnSuccess empty implementation.
func (*Noop) OnSuccess() {}

// OnFailure empty implementation.
func (*Noop) OnFailure() {}
