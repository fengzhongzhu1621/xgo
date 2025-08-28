package zipfile

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFromZip(t *testing.T) {
	// 保存解压缩出来的文件内容
	buf := bytes.NewBuffer(nil)

	// 测试解压缩文件到 buf 中，只读取 foo.txt 文件的内容
	err := ExtractFromZip("../../tests/testdata/test.zip", "**/foo.txt", buf)
	assert.Nil(t, err)
	t.Log("Content: " + buf.String())
}
