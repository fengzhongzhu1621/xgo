package compress

import (
	"errors"
	"testing"
)

func BenchmarkCheckNoopCompression(b *testing.B) {
	b.Run("check noop compression", func(b *testing.B) {
		bs := make([]byte, 1024)
		var result []byte
		for i := 0; i < b.N; i++ {
			result, _ = Compress(CompressTypeNoop, bs)
			bs, _ = Decompress(CompressTypeNoop, result)
		}
	})
	b.Run("not check noop compression", func(b *testing.B) {
		bs := make([]byte, 1024)
		var result []byte
		for i := 0; i < b.N; i++ {
			result, _ = oldCompress(CompressTypeNoop, bs)
			bs, _ = oldDecompress(CompressTypeNoop, result)
		}
	})
}

// oldCompress returns the compressed data, the data is compressed
// by a specific compressor.
func oldCompress(compressorType int, in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, nil
	}
	compressor := GetCompressor(compressorType)
	if compressor == nil {
		return nil, errors.New("compressor not registered")
	}
	return compressor.Compress(in)
}

// oldDecompress returns the decompressed data, the data is decompressed
// by a specific compressor.
func oldDecompress(compressorType int, in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, nil
	}
	compressor := GetCompressor(compressorType)
	if compressor == nil {
		return nil, errors.New("compressor not registered")
	}
	return compressor.Decompress(in)
}
