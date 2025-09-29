package codec

import "io"

// FramerBuilder defines how to build a framer. In general, each connection
// build a framer.
type FramerBuilder interface {
	New(io.Reader) Framer
}
