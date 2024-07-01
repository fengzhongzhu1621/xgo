package email

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin"
)

const (
	MaxLineLength      = 76                             // MaxLineLength is the maximum line length per RFC 2045
	defaultContentType = "text/plain; charset=us-ascii" // defaultContentType is the default Content-Type according to RFC 2045, section 5.2
)

// ErrMissingBoundary is returned when there is no boundary given for a multipart entity
var ErrMissingBoundary = errors.New("No boundary found for multipart entity")

// ErrMissingContentType is returned when there is no "Content-Type" header for a MIME entity
var ErrMissingContentType = errors.New("No Content-Type found for MIME entity")

// Email is the type used for email messages
type Email struct {
	ReplyTo     []string
	From        string   // 发件人
	To          []string // 收件人
	Bcc         []string
	Cc          []string
	Subject     string // 邮件标题
	Text        []byte // Plaintext message (optional) 文本格式的邮件
	HTML        []byte // Html message (optional)
	Sender      string // override From as SMTP envelope sender (optional)
	Headers     textproto.MIMEHeader
	Attachments []*Attachment // 附件
	ReadReceipt []string
}

// part is a copyable representation of a multipart.Part
type part struct {
	header textproto.MIMEHeader
	body   []byte
}

// Attach is used to attach content from an io.Reader to the email.
// Required parameters include an io.Reader, the desired filename for the attachment, and the Content-Type
// The function will return the created Attachment for reference, as well as nil for the error, if successful.
func (e *Email) Attach(r io.Reader, filename string, c string) (a *Attachment, err error) {
	// 读取附件内容
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, r); err != nil {
		return
	}
	at := &Attachment{
		Filename:    filename, // 附件名称
		ContentType: c,        // 附件的类型
		Header:      textproto.MIMEHeader{},
		Content:     buffer.Bytes(), // 附件内容
	}
	e.Attachments = append(e.Attachments, at)
	return at, nil
}

// AttachFile is used to attach content to the email.
// It attempts to open the file referenced by filename and, if successful, creates an Attachment.
// This Attachment is then appended to the slice of Email.Attachments.
// The function will then return the Attachment for reference, as well as nil for the error, if successful.
func (e *Email) AttachFile(filename string) (a *Attachment, err error) {
	// 打开附件
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	// 获得附件的类型
	ct := mime.TypeByExtension(filepath.Ext(filename))
	// 获得文件名，去掉前置路径
	basename := filepath.Base(filename)
	return e.Attach(f, basename, ct)
}

// 构造邮件的消息头
// msgHeaders merges the Email's various fields and custom headers together in a
// standards compliant way to create a MIMEHeader to be used in the resulting
// message. It does not alter e.Headers.
//
// "e"'s fields To, Cc, From, Subject will be used unless they are present in
// e.Headers. Unless set in e.Headers, "Date" will filled with the current time.
func (e *Email) msgHeaders() (textproto.MIMEHeader, error) {
	res := make(textproto.MIMEHeader, len(e.Headers)+6)
	if e.Headers != nil {
		for _, h := range []string{"Reply-To", "To", "Cc", "From", "Subject", "Date", "Message-Id", "MIME-Version"} {
			if v, ok := e.Headers[h]; ok {
				res[h] = v
			}
		}
	}
	// Set headers if there are values.
	if _, ok := res["Reply-To"]; !ok && len(e.ReplyTo) > 0 {
		res.Set("Reply-To", strings.Join(e.ReplyTo, ", "))
	}
	if _, ok := res["To"]; !ok && len(e.To) > 0 {
		res.Set("To", strings.Join(e.To, ", "))
	}
	if _, ok := res["Cc"]; !ok && len(e.Cc) > 0 {
		res.Set("Cc", strings.Join(e.Cc, ", "))
	}
	if _, ok := res["Subject"]; !ok && e.Subject != "" {
		res.Set("Subject", e.Subject)
	}
	if _, ok := res["Message-Id"]; !ok {
		id, err := generateMessageID()
		if err != nil {
			return nil, err
		}
		res.Set("Message-Id", id)
	}
	// Date and From are required headers.
	if _, ok := res["From"]; !ok {
		res.Set("From", e.From)
	}
	if _, ok := res["Date"]; !ok {
		res.Set("Date", time.Now().Format(time.RFC1123Z))
	}
	if _, ok := res["MIME-Version"]; !ok {
		res.Set("MIME-Version", "1.0")
	}
	for field, vals := range e.Headers {
		if _, ok := res[field]; !ok {
			res[field] = vals
		}
	}
	return res, nil
}

func writeMessage(buff io.Writer, msg []byte, multipart bool, mediaType string, w *multipart.Writer) error {
	if multipart {
		header := textproto.MIMEHeader{
			"Content-Type":              {mediaType + "; charset=UTF-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		}
		if _, err := w.CreatePart(header); err != nil {
			return err
		}
	}

	qp := quotedprintable.NewWriter(buff)
	// Write the text
	if _, err := qp.Write(msg); err != nil {
		return err
	}
	return qp.Close()
}

func (e *Email) categorizeAttachments() (htmlRelated, others []*Attachment) {
	for _, a := range e.Attachments {
		if a.HTMLRelated {
			htmlRelated = append(htmlRelated, a)
		} else {
			others = append(others, a)
		}
	}
	return
}

// Bytes converts the Email object to a []byte representation, including all needed MIMEHeaders, boundaries, etc.
func (e *Email) Bytes() ([]byte, error) {
	// TODO: better guess buffer size
	buff := bytes.NewBuffer(make([]byte, 0, 4096))

	headers, err := e.msgHeaders()
	if err != nil {
		return nil, err
	}

	htmlAttachments, otherAttachments := e.categorizeAttachments()
	if len(e.HTML) == 0 && len(htmlAttachments) > 0 {
		return nil, errors.New("there are HTML attachments, but no HTML body")
	}

	var (
		isMixed       = len(otherAttachments) > 0
		isAlternative = len(e.Text) > 0 && len(e.HTML) > 0
		isRelated     = len(e.HTML) > 0 && len(htmlAttachments) > 0
	)

	var w *multipart.Writer
	if isMixed || isAlternative || isRelated {
		w = multipart.NewWriter(buff)
	}
	switch {
	case isMixed:
		headers.Set("Content-Type", "multipart/mixed;\r\n boundary="+w.Boundary())
	case isAlternative:
		headers.Set("Content-Type", "multipart/alternative;\r\n boundary="+w.Boundary())
	case isRelated:
		headers.Set("Content-Type", "multipart/related;\r\n boundary="+w.Boundary())
	case len(e.HTML) > 0:
		headers.Set("Content-Type", "text/html; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	default:
		headers.Set("Content-Type", "text/plain; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	}
	headerToBytes(buff, headers)
	_, err = io.WriteString(buff, "\r\n")
	if err != nil {
		return nil, err
	}

	// Check to see if there is a Text or HTML field
	if len(e.Text) > 0 || len(e.HTML) > 0 {
		var subWriter *multipart.Writer

		if isMixed && isAlternative {
			// Create the multipart alternative part
			subWriter = multipart.NewWriter(buff)
			header := textproto.MIMEHeader{
				"Content-Type": {"multipart/alternative;\r\n boundary=" + subWriter.Boundary()},
			}
			if _, err := w.CreatePart(header); err != nil {
				return nil, err
			}
		} else {
			subWriter = w
		}
		// Create the body sections
		if len(e.Text) > 0 {
			// Write the text
			if err := writeMessage(buff, e.Text, isMixed || isAlternative, "text/plain", subWriter); err != nil {
				return nil, err
			}
		}
		if len(e.HTML) > 0 {
			messageWriter := subWriter
			var relatedWriter *multipart.Writer
			if (isMixed || isAlternative) && len(htmlAttachments) > 0 {
				relatedWriter = multipart.NewWriter(buff)
				header := textproto.MIMEHeader{
					"Content-Type": {"multipart/related;\r\n boundary=" + relatedWriter.Boundary()},
				}
				if _, err := subWriter.CreatePart(header); err != nil {
					return nil, err
				}

				messageWriter = relatedWriter
			} else if isRelated && len(htmlAttachments) > 0 {
				relatedWriter = w
				messageWriter = w
			}
			// Write the HTML
			if err := writeMessage(buff, e.HTML, isMixed || isAlternative || isRelated, "text/html", messageWriter); err != nil {
				return nil, err
			}
			if len(htmlAttachments) > 0 {
				for _, a := range htmlAttachments {
					a.setDefaultHeaders()
					ap, err := relatedWriter.CreatePart(a.Header)
					if err != nil {
						return nil, err
					}
					// Write the base64Wrapped content to the part
					base64Wrap(ap, a.Content)
				}

				if isMixed || isAlternative {
					relatedWriter.Close()
				}
			}
		}
		if isMixed && isAlternative {
			if err := subWriter.Close(); err != nil {
				return nil, err
			}
		}
	}
	// Create attachment part, if necessary
	for _, a := range otherAttachments {
		a.setDefaultHeaders()
		ap, err := w.CreatePart(a.Header)
		if err != nil {
			return nil, err
		}
		// Write the base64Wrapped content to the part
		base64Wrap(ap, a.Content)
	}
	if isMixed || isAlternative || isRelated {
		if err := w.Close(); err != nil {
			return nil, err
		}
	}
	return buff.Bytes(), nil
}

// Send an email using the given host and SMTP auth (optional), returns any error thrown by smtp.SendMail
// This function merges the To, Cc, and Bcc fields and calls the smtp.SendMail function using the Email.Bytes() output as the message
func (e *Email) Send(addr string, a smtp.Auth) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
	to = append(append(append(to, e.To...), e.Cc...), e.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if e.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := e.parseSender()
	if err != nil {
		return err
	}
	raw, err := e.Bytes()
	if err != nil {
		return err
	}
	return smtp.SendMail(addr, a, sender, to, raw)
}

// Select and parse an SMTP envelope sender address.  Choose Email.Sender if set, or fallback to Email.From.
func (e *Email) parseSender() (string, error) {
	if e.Sender != "" {
		sender, err := mail.ParseAddress(e.Sender)
		if err != nil {
			return "", err
		}
		return sender.Address, nil
	} else {
		from, err := mail.ParseAddress(e.From)
		if err != nil {
			return "", err
		}
		return from.Address, nil
	}
}

// SendWithTLS sends an email over tls with an optional TLS config.
//
// The TLS Config is helpful if you need to connect to a host that is used an untrusted
// certificate.
func (e *Email) SendWithTLS(addr string, a smtp.Auth, t *tls.Config) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
	to = append(append(append(to, e.To...), e.Cc...), e.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if e.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := e.parseSender()
	if err != nil {
		return err
	}
	raw, err := e.Bytes()
	if err != nil {
		return err
	}

	conn, err := tls.Dial("tcp", addr, t)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, t.ServerName)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(sender); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// SendWithStartTLS sends an email over TLS using STARTTLS with an optional TLS config.
//
// The TLS Config is helpful if you need to connect to a host that is used an untrusted
// certificate.
func (e *Email) SendWithStartTLS(addr string, a smtp.Auth, t *tls.Config) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(e.To)+len(e.Cc)+len(e.Bcc))
	to = append(append(append(to, e.To...), e.Cc...), e.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if e.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := e.parseSender()
	if err != nil {
		return err
	}
	raw, err := e.Bytes()
	if err != nil {
		return err
	}

	// Taken from the standard library
	// https://github.com/golang/go/blob/master/src/net/smtp/smtp.go#L328
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("localhost"); err != nil {
		return err
	}
	// Use TLS if available
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(t); err != nil {
			return err
		}
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(sender); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(raw)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// NewEmail creates an Email, and returns the pointer to it.
func NewEmail() *Email {
	return &Email{Headers: textproto.MIMEHeader{}}
}

// NewEmailFromReader reads a stream of bytes from an io.Reader, r,
// and returns an email struct containing the parsed data.
// This function expects the data in RFC 5322 format.
func NewEmailFromReader(r io.Reader) (*Email, error) {
	e := NewEmail()
	s := &buildin.TrimReader{Rd: r}
	tp := textproto.NewReader(bufio.NewReader(s))
	// Parse the main headers
	hdrs, err := tp.ReadMIMEHeader()
	if err != nil {
		return e, err
	}
	// Set the subject, to, cc, bcc, and from
	for h, v := range hdrs {
		switch h {
		case "Subject":
			e.Subject = v[0]
			subj, err := (&mime.WordDecoder{}).DecodeHeader(e.Subject)
			if err == nil && len(subj) > 0 {
				e.Subject = subj
			}
			delete(hdrs, h)
		case "To":
			e.To = handleAddressList(v)
			delete(hdrs, h)
		case "Cc":
			e.Cc = handleAddressList(v)
			delete(hdrs, h)
		case "Bcc":
			e.Bcc = handleAddressList(v)
			delete(hdrs, h)
		case "Reply-To":
			e.ReplyTo = handleAddressList(v)
			delete(hdrs, h)
		case "From":
			e.From = v[0]
			fr, err := (&mime.WordDecoder{}).DecodeHeader(e.From)
			if err == nil && len(fr) > 0 {
				e.From = fr
			}
			delete(hdrs, h)
		}
	}
	e.Headers = hdrs
	body := tp.R
	// Recursively parse the MIME parts
	ps, err := parseMIMEParts(e.Headers, body)
	if err != nil {
		return e, err
	}
	for _, p := range ps {
		if ct := p.header.Get("Content-Type"); ct == "" {
			return e, ErrMissingContentType
		}
		ct, _, err := mime.ParseMediaType(p.header.Get("Content-Type"))
		if err != nil {
			return e, err
		}
		// Check if part is an attachment based on the existence of the Content-Disposition header with a value of "attachment".
		if cd := p.header.Get("Content-Disposition"); cd != "" {
			cd, params, err := mime.ParseMediaType(p.header.Get("Content-Disposition"))
			if err != nil {
				return e, err
			}
			filename, filenameDefined := params["filename"]
			if cd == "attachment" || (cd == "inline" && filenameDefined) {
				_, err = e.Attach(bytes.NewReader(p.body), filename, ct)
				if err != nil {
					return e, err
				}
				continue
			}
		}
		switch {
		case ct == "text/plain":
			e.Text = p.body
		case ct == "text/html":
			e.HTML = p.body
		}
	}
	return e, nil
}
