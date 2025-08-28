package publisher

import "errors"

// ErrOutputInNoPublisherHandler happens when a handler func returned some messages in a no-publisher handler.
// todo: maybe change the handler func signature in no-publisher handler so that there's no possibility for this
var ErrOutputInNoPublisherHandler = errors.New(
	"returned output messages in a handler without publisher",
)
