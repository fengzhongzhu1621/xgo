package throttle

// Throttler defines the throttle interface.
type Throttler interface {
	Allow() bool
	OnSuccess()
	OnFailure()
}
