package compress

import (
	"errors"
)

// ICompressor is body compress and decompress interface.
type ICompressor interface {
	Compress(in []byte) (out []byte, err error)
	Decompress(in []byte) (out []byte, err error)
}

// CompressType is the mode of body compress or decompress.
const (
	CompressTypeNoop = iota
	CompressTypeGzip
	CompressTypeSnappy
	CompressTypeZlib
	CompressTypeStreamSnappy
	CompressTypeBlockSnappy
)

var compressors = make(map[int]ICompressor)

// RegisterCompressor register a specific compressor, which will
// be called by init function defined in third package.
func RegisterCompressor(compressType int, s ICompressor) {
	compressors[compressType] = s
}

// GetCompressor returns a specific compressor by type.
func GetCompressor(compressType int) ICompressor {
	return compressors[compressType]
}

// Compress returns the compressed data, the data is compressed
// by a specific compressor.
func Compress(compressorType int, in []byte) ([]byte, error) {
	// Explicitly check for noop to avoid accessing the map.
	if compressorType == CompressTypeNoop {
		return in, nil
	}
	if len(in) == 0 {
		return nil, nil
	}
	compressor := GetCompressor(compressorType)
	if compressor == nil {
		return nil, errors.New("compressor not registered")
	}
	return compressor.Compress(in)
}

// Decompress returns the decompressed data, the data is decompressed
// by a specific compressor.
func Decompress(compressorType int, in []byte) ([]byte, error) {
	// Explicitly check for noop to avoid accessing the map.
	if compressorType == CompressTypeNoop {
		return in, nil
	}
	if len(in) == 0 {
		return nil, nil
	}
	compressor := GetCompressor(compressorType)
	if compressor == nil {
		return nil, errors.New("compressor not registered")
	}
	return compressor.Decompress(in)
}
