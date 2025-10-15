package codec

import "io"

// IFramerBuilder defines how to build a framer. In general, each connection
// build a framer.
type IFramerBuilder interface {
	New(io.Reader) IFramer
}
