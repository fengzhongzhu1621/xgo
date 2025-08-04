package email

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"sync"
	"syscall"
	"time"
)

const maxFails = 4

var (
	ErrClosed  = errors.New("pool closed")
	ErrTimeout = errors.New("timed out")
)

type client struct {
	*smtp.Client
	failCount int
}

// go1.1 didn't have this method
func (c *client) Close() error {
	return c.Text.Close()
}

type timestampedErr struct {
	err error
	ts  time.Time
}

type Pool struct {
	addr          string
	auth          smtp.Auth
	max           int // 客户端的最大数量
	created       int // 已连接的客户端的数量
	clients       chan *client
	rebuild       chan struct{}
	mut           *sync.Mutex
	lastBuildErr  *timestampedErr
	closing       chan struct{}
	tlsConfig     *tls.Config
	helloHostname string
}

func NewPool(
	address string,
	count int,
	auth smtp.Auth,
	opt_tlsConfig ...*tls.Config,
) (pool *Pool, err error) {
	pool = &Pool{
		addr:    address,
		auth:    auth,
		max:     count,
		clients: make(chan *client, count), // 客户端的数量
		rebuild: make(chan struct{}),
		closing: make(chan struct{}),
		mut:     &sync.Mutex{},
	}
	if len(opt_tlsConfig) == 1 {
		pool.tlsConfig = opt_tlsConfig[0]
	} else if host, _, e := net.SplitHostPort(address); e != nil {
		return nil, e
	} else {
		pool.tlsConfig = &tls.Config{ServerName: host}
	}
	return
}

// SetHelloHostname optionally sets the hostname that the Go smtp.Client will
// use when doing a HELLO with the upstream SMTP server. By default, Go uses
// "localhost" which may not be accepted by certain SMTP servers that demand
// an FQDN.
func (p *Pool) SetHelloHostname(h string) {
	p.helloHostname = h
}

func (p *Pool) get(timeout time.Duration) *client {
	// 获得可用的客户端
	select {
	case c := <-p.clients:
		return c
	default:
	}
	// 如果没有可用的客户端，则创建一个新的客户端
	if p.created < p.max {
		p.makeOne()
	}

	// 创建超时事件
	var deadline <-chan time.Time
	if timeout >= 0 {
		deadline = time.After(timeout)
	}

	for {
		select {
		case c := <-p.clients:
			return c
		case <-p.rebuild:
			p.makeOne()
		case <-deadline:
			// 超时获取客户端失败
			return nil
		case <-p.closing:
			return nil
		}
	}
}

// 将连接重新放到连接池
func (p *Pool) replace(c *client) {
	p.clients <- c
}

// 将可用的客户端加一，计数增加了，但是实际上创建连接失败
func (p *Pool) inc() bool {
	if p.created >= p.max {
		return false
	}

	p.mut.Lock()
	defer p.mut.Unlock()

	if p.created >= p.max {
		return false
	}
	p.created++
	return true
}

func (p *Pool) dec() {
	p.mut.Lock()
	p.created--
	p.mut.Unlock()
	// 发送重建连接信号
	select {
	case p.rebuild <- struct{}{}:
	default:
	}
}

// 创建一个新的客户端
func (p *Pool) makeOne() {
	go func() {
		// 将可用的客户端加一
		if p.inc() {
			if c, err := p.build(); err == nil {
				p.clients <- c
			} else {
				p.lastBuildErr = &timestampedErr{err, time.Now()}
				// 回滚次数
				p.dec()
			}
		}
	}()
}

// 创建连接
func (p *Pool) build() (*client, error) {
	cl, err := smtp.Dial(p.addr)
	if err != nil {
		return nil, err
	}

	// Is there a custom hostname for doing a HELLO with the SMTP server?
	if p.helloHostname != "" {
		cl.Hello(p.helloHostname)
	}

	c := &client{cl, 0}

	if _, err := startTLS(c, p.tlsConfig); err != nil {
		c.Close()
		return nil, err
	}

	if p.auth != nil {
		if _, err := addAuth(c, p.auth); err != nil {
			c.Close()
			return nil, err
		}
	}

	return c, nil
}

// 发送之后的处理
func (p *Pool) maybeReplace(err error, c *client) {
	// 邮件发送成功
	if err == nil {
		c.failCount = 0
		// 将连接重新放到连接池
		p.replace(c)
		return
	}
	// 记录发送失败次数
	c.failCount++
	if c.failCount >= maxFails {
		goto shutdown
	}

	if !shouldReuse(err) {
		goto shutdown
	}

	if err := c.Reset(); err != nil {
		goto shutdown
	}
	// 将连接重新放到连接池
	p.replace(c)
	return

shutdown:
	p.dec()
	c.Close()
}

func (p *Pool) failedToGet(startTime time.Time) error {
	select {
	case <-p.closing:
		return ErrClosed
	default:
	}

	if p.lastBuildErr != nil && startTime.Before(p.lastBuildErr.ts) {
		return p.lastBuildErr.err
	}

	return ErrTimeout
}

// Send sends an email via a connection pulled from the Pool. The timeout may
// be <0 to indicate no timeout. Otherwise reaching the timeout will produce
// and error building a connection that occurred while we were waiting, or
// otherwise ErrTimeout.
func (p *Pool) Send(e *Email, timeout time.Duration) (err error) {
	start := time.Now()
	// 获得客户端连接，如果没有则创建一个新的连接
	c := p.get(timeout)
	if c == nil {
		// 返回自定义异常
		return p.failedToGet(start)
	}

	defer func() {
		// 如果发送失败，则将连接重新放回连接池
		p.maybeReplace(err, c)
	}()

	recipients, err := addressLists(e.To, e.Cc, e.Bcc)
	if err != nil {
		return
	}

	msg, err := e.Bytes()
	if err != nil {
		return
	}

	from, err := emailOnly(e.From)
	if err != nil {
		return
	}
	if err = c.Mail(from); err != nil {
		return
	}

	for _, recip := range recipients {
		if err = c.Rcpt(recip); err != nil {
			return
		}
	}

	w, err := c.Data()
	if err != nil {
		return
	}
	if _, err = w.Write(msg); err != nil {
		return
	}

	err = w.Close()

	return
}

// Close immediately changes the pool's state so no new connections will be
// created, then gets and closes the existing ones as they become available.
func (p *Pool) Close() {
	close(p.closing)

	for p.created > 0 {
		c := <-p.clients
		c.Quit()
		p.dec()
	}
}

func shouldReuse(err error) bool {
	// certainly not perfect, but might be close:
	//  - EOF: clearly, the connection went down
	//  - textproto.Errors were valid SMTP over a valid connection,
	//    but resulted from an SMTP error response
	//  - textproto.ProtocolErrors result from connections going down,
	//    invalid SMTP, that sort of thing
	//  - syscall.Errno is probably down connection/bad pipe, but
	//    passed straight through by textproto instead of becoming a
	//    ProtocolError
	//  - if we don't recognize the error, don't reuse the connection
	// A false positive will probably fail on the Reset(), and even if
	// not will eventually hit maxFails.
	// A false negative will knock over (and trigger replacement of) a
	// conn that might have still worked.
	if err == io.EOF {
		return false
	}
	switch err.(type) {
	case *textproto.Error:
		return true
	case *textproto.ProtocolError, textproto.ProtocolError:
		return false
	case syscall.Errno:
		return false
	default:
		return false
	}
}

func startTLS(c *client, t *tls.Config) (bool, error) {
	if ok, _ := c.Extension("STARTTLS"); !ok {
		return false, nil
	}

	if err := c.StartTLS(t); err != nil {
		return false, err
	}

	return true, nil
}

func addAuth(c *client, auth smtp.Auth) (bool, error) {
	if ok, _ := c.Extension("AUTH"); !ok {
		return false, nil
	}

	if err := c.Auth(auth); err != nil {
		return false, err
	}

	return true, nil
}
