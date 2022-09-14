package publisher

import "errors"

var (
	// ErrOutputInNoPublisherHandler happens when a handler func returned some messages in a no-publisher handler.
	// todo: maybe change the handler func signature in no-publisher handler so that there's no possibility for this
	ErrOutputInNoPublisherHandler = errors.New("returned output messages in a handler without publisher")
)
