package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

func init() {
	RegisterCompressor(CompressTypeGzip, &GzipCompress{})
}

var _ ICompressor = (*GzipCompress)(nil)

// GzipCompress is gzip compressor.
// GzipCompress 是gzip压缩器实现
type GzipCompress struct {
	readerPool sync.Pool // 用于复用 gzip.Reader 对象的对象池，减少内存分配开销
	writerPool sync.Pool // 用于复用 gzip.Writer 对象的对象池，减少内存分配开销
}

// Compress returns binary data compressed by gzip.
// Compress 使用gzip压缩二进制数据并返回结果
func (c *GzipCompress) Compress(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return in, nil
	}

	buffer := &bytes.Buffer{} // 创建字节缓冲区存储压缩结果
	// 从对象池获取或创建gzip写入器
	z, ok := c.writerPool.Get().(*gzip.Writer)
	if !ok {
		z = gzip.NewWriter(buffer) // 创建新的gzip写入器
	} else {
		z.Reset(buffer) // 重置已有的写入器指向新的缓冲区
	}
	defer c.writerPool.Put(z) // 函数结束时将写入器放回对象池

	// 写入数据到gzip流
	if _, err := z.Write(in); err != nil {
		return nil, err
	}
	// 关闭gzip流，完成压缩
	if err := z.Close(); err != nil {
		return nil, err
	}

	// 返回压缩后的字节数据
	return buffer.Bytes(), nil
}

// Decompress returns binary data decompressed by gzip.
// Decompress 使用gzip解压缩二进制数据并返回结果
func (c *GzipCompress) Decompress(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return in, nil
	}

	// 从对象池获取gzip读取器
	z, ok := c.readerPool.Get().(*gzip.Reader)

	// 函数结束时将读取器放回对象池
	defer func() {
		if z != nil {
			c.readerPool.Put(z)
		}
	}()

	// 如果对象池中没有可用的读取器，创建新的
	br := bytes.NewReader(in) // 创建字节读取器
	if !ok {
		gr, err := gzip.NewReader(br)
		if err != nil {
			return nil, err
		}
		z = gr
	} else {
		// 重置已有的读取器指向新的数据源
		if err := z.Reset(br); err != nil {
			return nil, err
		}
	}
	// 读取所有解压缩的数据
	out, err := io.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return out, nil
}
