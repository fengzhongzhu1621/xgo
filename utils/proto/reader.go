package proto

import (
	"bufio"
	"fmt"
	"io"

	"xgo/utils/bytesconv"
)

// redis resp protocol data type.
const (
	ErrorReply  = '-'
	StatusReply = '+'
	IntReply    = ':'
	StringReply = '$'
	ArrayReply  = '*'
)

//------------------------------------------------------------------------------

const Nil = RedisError("redis: nil") // nolint:errname

type RedisError string

func (e RedisError) Error() string { return string(e) }

func (RedisError) RedisError() {}

//------------------------------------------------------------------------------

type MultiBulkParse func(*Reader, int64) (interface{}, error)

type Reader struct {
	rd   *bufio.Reader
	_buf []byte			// 读缓存大小
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd:   bufio.NewReader(rd),
		_buf: make([]byte, 64),		// 读缓存大小，默认64字节
	}
}

// buf中可被读的长度（缓存数据的大小）
func (r *Reader) Buffered() int {
	return r.rd.Buffered()
}

// 只是“窥探”一下 Reader 中没有读取的 n 个字节。好比栈数据结构中的取栈顶元素，但不出栈。
// 返回的 []byte 只是 buffer 中的引用，在下次IO操作后会无效，可见该方法（以及ReadSlice这样的，返回buffer引用的方法）对多 goroutine 是不安全的，也就是在多并发环境下，不能依赖其结果。
func (r *Reader) Peek(n int) ([]byte, error) {
	return r.rd.Peek(n)
}

// 重置缓存，丢弃未被处理的数据
func (r *Reader) Reset(rd io.Reader) {
	r.rd.Reset(rd)
}

// 读取一行数据，并去掉换行符
func (r *Reader) ReadLine() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	if isNilReply(line) {
		return nil, Nil
	}
	return line, nil
}

// readLine that returns an error if:
//   - there is a pending read error;
//   - or line does not end with \r\n.
func (r *Reader) readLine() ([]byte, error) {
	// 读取换行符前的字节
	b, err := r.rd.ReadSlice('\n')
	if err != nil {
		if err != bufio.ErrBufferFull {
			return nil, err
		}

		full := make([]byte, len(b))
		copy(full, b)
		// 读取换行符
		b, err = r.rd.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		full = append(full, b...) //nolint:makezero
		b = full
	}
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("redis: invalid reply: %q", b)
	}
	// 去掉回车换行
	return b[:len(b)-2], nil
}

func (r *Reader) ReadReply(m MultiBulkParse) (interface{}, error) {
	line, err := r.ReadLine()
	fmt.Println("line = ", line)
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case StatusReply:
		return string(line[1:]), nil
	case IntReply:
		return bytesconv.ParseInt(line[1:], 10, 64)
	case StringReply:
		return r.readStringReply(line)
	case ArrayReply:
		// 数组的长度
		n, err := parseArrayLen(line)
		if err != nil {
			return nil, err
		}
		if m == nil {
			err := fmt.Errorf("redis: got %.100q, but multi bulk parser is nil", line)
			return nil, err
		}
		// 解析后面的数组内容
		return m(r, n)
	}
	return nil, fmt.Errorf("redis: can't parse %.100q", line)
}

func (r *Reader) ReadIntReply() (int64, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	switch line[0] {
	case ErrorReply:
		return 0, ParseErrorReply(line)
	case IntReply:
		return bytesconv.ParseInt(line[1:], 10, 64)
	default:
		return 0, fmt.Errorf("redis: can't parse int reply: %.100q", line)
	}
}

func (r *Reader) ReadString() (string, error) {
	line, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	switch line[0] {
	case ErrorReply:
		return "", ParseErrorReply(line)
	case StringReply:
		return r.readStringReply(line)
	case StatusReply:
		return string(line[1:]), nil
	case IntReply:
		return string(line[1:]), nil
	default:
		return "", fmt.Errorf("redis: can't parse reply=%.100q reading string", line)
	}
}

func (r *Reader) readStringReply(line []byte) (string, error) {
	if isNilReply(line) {
		return "", Nil
	}

	replyLen, err := bytesconv.Atoi(line[1:])
	if err != nil {
		return "", err
	}

	b := make([]byte, replyLen+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil {
		return "", err
	}

	return bytesconv.BytesToString(b[:replyLen]), nil
}

func (r *Reader) ReadArrayReply(m MultiBulkParse) (interface{}, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case ArrayReply:
		n, err := parseArrayLen(line)
		if err != nil {
			return nil, err
		}
		return m(r, n)
	default:
		return nil, fmt.Errorf("redis: can't parse array reply: %.100q", line)
	}
}

func (r *Reader) ReadArrayLen() (int, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	switch line[0] {
	case ErrorReply:
		return 0, ParseErrorReply(line)
	case ArrayReply:
		n, err := parseArrayLen(line)
		if err != nil {
			return 0, err
		}
		return int(n), nil
	default:
		return 0, fmt.Errorf("redis: can't parse array reply: %.100q", line)
	}
}

func (r *Reader) ReadScanReply() ([]string, uint64, error) {
	n, err := r.ReadArrayLen()
	if err != nil {
		return nil, 0, err
	}
	if n != 2 {
		return nil, 0, fmt.Errorf("redis: got %d elements in scan reply, expected 2", n)
	}

	cursor, err := r.ReadUint()
	if err != nil {
		return nil, 0, err
	}

	n, err = r.ReadArrayLen()
	if err != nil {
		return nil, 0, err
	}

	keys := make([]string, n)

	for i := 0; i < n; i++ {
		key, err := r.ReadString()
		if err != nil {
			return nil, 0, err
		}
		keys[i] = key
	}

	return keys, cursor, err
}

func (r *Reader) ReadInt() (int64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return bytesconv.ParseInt(b, 10, 64)
}

func (r *Reader) ReadUint() (uint64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return bytesconv.ParseUint(b, 10, 64)
}

func (r *Reader) ReadFloatReply() (float64, error) {
	b, err := r.readTmpBytesReply()
	if err != nil {
		return 0, err
	}
	return bytesconv.ParseFloat(b, 64)
}

func (r *Reader) readTmpBytesReply() ([]byte, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	switch line[0] {
	case ErrorReply:
		return nil, ParseErrorReply(line)
	case StringReply:
		return r._readTmpBytesReply(line)
	case StatusReply:
		return line[1:], nil
	default:
		return nil, fmt.Errorf("redis: can't parse string reply: %.100q", line)
	}
}

func (r *Reader) _readTmpBytesReply(line []byte) ([]byte, error) {
	if isNilReply(line) {
		return nil, Nil
	}

	replyLen, err := bytesconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	buf := r.buf(replyLen + 2)
	_, err = io.ReadFull(r.rd, buf)
	if err != nil {
		return nil, err
	}

	return buf[:replyLen], nil
}

func (r *Reader) buf(n int) []byte {
	if n <= cap(r._buf) {
		return r._buf[:n]
	}
	d := n - cap(r._buf)
	r._buf = append(r._buf, make([]byte, d)...)
	return r._buf
}

func isNilReply(b []byte) bool {
	return len(b) == 3 &&
		(b[0] == StringReply || b[0] == ArrayReply) &&
		b[1] == '-' && b[2] == '1'
}

func ParseErrorReply(line []byte) error {
	return RedisError(string(line[1:]))
}

func parseArrayLen(line []byte) (int64, error) {
	if isNilReply(line) {
		return 0, Nil
	}
	return bytesconv.ParseInt(line[1:], 10, 64)
}
