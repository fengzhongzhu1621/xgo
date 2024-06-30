package encoding

import (
	"sync"

	"github.com/fengzhongzhu1621/xgo"
)

// Encoder encodes the contents of v into a byte representation.
// It's primarily used for encoding a map[string]interface{} into a file format.
type Encoder interface {
	// Encode 将字符串字典转换为字节数组
	Encode(v map[string]interface{}) ([]byte, error)
}

const (
	// ErrEncoderNotFound is returned when there is no encoder registered for a format.
	ErrEncoderNotFound = xgo.EncodingError("encoder not found for this format")

	// ErrEncoderFormatAlreadyRegistered is returned when an encoder is already registered for a format.
	ErrEncoderFormatAlreadyRegistered = xgo.EncodingError("encoder already registered for this format")
)

// EncoderRegistry can choose an appropriate Encoder based on the provided format
// 编码器仓库 .
type EncoderRegistry struct {
	// 存放编码器
	encoders map[string]Encoder

	mu sync.RWMutex
}

// NewEncoderRegistry returns a new, initialized EncoderRegistry.
func NewEncoderRegistry() *EncoderRegistry {
	return &EncoderRegistry{
		encoders: make(map[string]Encoder),
	}
}

// RegisterEncoder registers an Encoder for a format.
// Registering an Encoder for an already existing format is not supported.
func (e *EncoderRegistry) RegisterEncoder(format string, enc Encoder) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 判断编码器是否已经注册，如果已经注册，返回错误
	if _, ok := e.encoders[format]; ok {
		return ErrEncoderFormatAlreadyRegistered
	}

	e.encoders[format] = enc

	return nil
}

// Encode 根据编码器的名称选择一个编码器对字典进行编码.
func (e *EncoderRegistry) Encode(format string, v map[string]interface{}) ([]byte, error) {
	e.mu.RLock()
	encoder, ok := e.encoders[format]
	e.mu.RUnlock()

	// 判断编码器是否存在
	if !ok {
		return nil, ErrEncoderNotFound
	}

	return encoder.Encode(v)
}
