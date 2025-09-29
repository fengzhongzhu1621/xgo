package attachment

import "io"

var _ io.Reader = (*NoopAttachment)(nil)

// NoopAttachment is an empty attachment.
type NoopAttachment struct{}

// Read implements the io.Reader interface, which always returns (0, io.EOF)
func (a NoopAttachment) Read(_ []byte) (n int, err error) {
	return 0, io.EOF
}
