package file

import (
	"os"
	"testing"

	"github.com/fengzhongzhu1621/xgo/crypto/randutils"
	"github.com/stretchr/testify/assert"
)

func TestSanitizedName(t *testing.T) {
	tests := []struct {
		orig   string
		expect string
	}{
		{"", "."},
		{"//../foo", "foo"},
		{"/../../", ""},
		{"/hello/world/..", "hello"},
		{"/..", ""},
		{"/foo/..", ""},
		{"/-/foo", "-/foo"},
	}
	for _, v := range tests {
		res := SanitizedName(v.orig)
		if res != v.expect {
			t.Fatalf("Clean path(%v) expect(%v) but got(%v)", v.orig, v.expect, res)
		}
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name string
		f    func()
	}{
		{
			name: "file exist",
			f: func() {
				// 创建临时文件
				f, err := os.CreateTemp("", "test")
				assert.NoError(t, err)
				defer os.Remove(f.Name())
				// 判读文件是否存在
				exist, err := FileOrDirExists(f.Name())
				assert.NoError(t, err)
				assert.True(t, exist)
			},
		},
		{
			name: "file not exist",
			f: func() {
				// 创建一个随机文件名，确保文件不存在
				fileName := randutils.RandomString(10)
				exist, err := FileOrDirExists(fileName)
				assert.NoError(t, err)
				assert.False(t, exist)
			},
		},
		{
			name: "is dir",
			f: func() {
				// 创建临时目录
				dirPath, err := os.MkdirTemp("", "example")
				assert.NoError(t, err)
				defer os.RemoveAll(dirPath) // clean up
				// 判断文件夹是否存在
				exist, err := FileOrDirExists(dirPath)
				assert.NoError(t, err)
				assert.True(t, exist)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f()
		})
	}
}
