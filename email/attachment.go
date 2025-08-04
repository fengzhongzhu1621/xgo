package email

import (
	"fmt"
	"net/textproto"
)

// Attachment is a struct representing an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	Filename    string // 附件名称
	ContentType string // 附件类型
	Header      textproto.MIMEHeader
	Content     []byte // 附件内容
	HTMLRelated bool
}

// 设置附件的header
func (at *Attachment) setDefaultHeaders() {
	// 附件的默认类型
	contentType := "application/octet-stream"
	if len(at.ContentType) > 0 {
		contentType = at.ContentType
	}
	at.Header.Set("Content-Type", contentType)

	if len(at.Header.Get("Content-Disposition")) == 0 {
		disposition := "attachment"
		if at.HTMLRelated {
			disposition = "inline"
		}
		at.Header.Set(
			"Content-Disposition",
			fmt.Sprintf("%s;\r\n filename=\"%s\"", disposition, at.Filename),
		)
	}
	if len(at.Header.Get("Content-ID")) == 0 {
		at.Header.Set("Content-ID", fmt.Sprintf("<%s>", at.Filename))
	}
	if len(at.Header.Get("Content-Transfer-Encoding")) == 0 {
		at.Header.Set("Content-Transfer-Encoding", "base64")
	}
}
