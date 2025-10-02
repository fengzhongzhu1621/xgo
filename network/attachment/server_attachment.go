package attachment

import (
	"io"

	"github.com/fengzhongzhu1621/xgo/codec"
)

// Attachment stores the attachment of tRPC requests/responses.
type ServerAttachment struct {
	attachment *Attachment
}

// Request returns Request Attachment.
func (a *ServerAttachment) Request() io.Reader {
	return a.attachment.Request
}

// SetResponse sets Response attachment.
func (a *ServerAttachment) SetResponse(attachment io.Reader) {
	a.attachment.Response = attachment
}

// GetServerAttachment returns Attachment from msg.
// If there is no Attachment in the msg, an empty attachment bound to the msg will be returned.
func GetServerAttachment(msg codec.IMsg) *ServerAttachment {
	cm := msg.CommonMeta()
	if cm == nil {
		cm = make(codec.CommonMeta)
		msg.WithCommonMeta(cm)
	}
	a, _ := cm[ServerAttachmentKey{}]
	if a == nil {
		// 附件不存在，则创建一个新的附件
		// 这里的NoopAttachment{}是一个空的附件，它实现了io.Reader接口，但每次读取时都会返回0和io.EOF错误。
		a = &Attachment{Request: NoopAttachment{}, Response: NoopAttachment{}}
		cm[ServerAttachmentKey{}] = a
	}

	return &ServerAttachment{attachment: a.(*Attachment)}
}
