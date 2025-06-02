package brotli

import "testing"

func TestCompress(t *testing.T) {
	original := []byte("This is some sample text to compress and decompress using Brotli.")

	// 压缩
	compressed, err := Compress(original)
	if err != nil {
		panic(err)
	}

	// 解压缩
	decompressed, err := Decompress(compressed)
	if err != nil {
		panic(err)
	}

	// 验证
	if string(original) != string(decompressed) {
		panic("Decompressed data doesn't match original")
	}
}
