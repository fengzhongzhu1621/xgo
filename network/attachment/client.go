package attachment

import (
	"io"
)

// ClientAttachment stores the Attachment of tRPC requests/responses.
type ClientAttachment struct {
	attachment Attachment
}

// NewAttachment returns a new Attachment whose response Attachment is a NoopAttachment.
func NewAttachment(request io.Reader) *ClientAttachment {
	return &ClientAttachment{attachment: Attachment{Request: request, Response: NoopAttachment{}}}
}

// Response returns Response Attachment.
func (a *ClientAttachment) Response() io.Reader {
	return a.attachment.Response
}
