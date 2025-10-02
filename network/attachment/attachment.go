package attachment

import (
	"bytes"
	"io"

	"github.com/fengzhongzhu1621/xgo/codec"
)

// ClientAttachmentKey is the key of client's Attachment.
type ClientAttachmentKey struct{}

// ServerAttachmentKey is the key of server's Attachment.
type ServerAttachmentKey struct{}

// Attachment stores the attachment in tRPC requests/responses.
type Attachment struct {
	Request  io.Reader
	Response io.Reader
}

// ClientRequestAttachment returns client's Request Attachment from msg.
func ClientRequestAttachment(msg codec.IMsg) (io.Reader, bool) {
	// 从msg中获取ClientAttachmentKey{}对应的值，并将其转换为*Attachment类型。
	if a, _ := msg.CommonMeta()[ClientAttachmentKey{}].(*Attachment); a != nil {
		return a.Request, true
	}
	return nil, false
}

// SetClientResponseAttachment sets client's Response attachment to msg.
// If the message does not contain client.Attachment,
// which means that the user has explicitly ignored the att returned by the server.
// For performance reasons, there is no need to set the response attachment into msg.
func SetClientResponseAttachment(msg codec.IMsg, attachment []byte) {
	if a, _ := msg.CommonMeta()[ClientAttachmentKey{}].(*Attachment); a != nil {
		// 存在则覆盖
		a.Response = bytes.NewReader(attachment)
	}
}

// ServerResponseAttachment returns server's Response Attachment from msg.
func ServerResponseAttachment(msg codec.IMsg) (io.Reader, bool) {
	if a, _ := msg.CommonMeta()[ServerAttachmentKey{}].(*Attachment); a != nil {
		return a.Response, true
	}
	return nil, false
}

// SetServerRequestAttachment sets server's Request Attachment to msg.
func SetServerRequestAttachment(m codec.IMsg, attachment []byte) {
	cm := m.CommonMeta()
	if cm == nil {
		// 没有则创建一个新的
		cm = make(codec.CommonMeta)
		m.WithCommonMeta(cm)
	}
	cm[ServerAttachmentKey{}] = &Attachment{Request: bytes.NewReader(attachment), Response: NoopAttachment{}}
}
