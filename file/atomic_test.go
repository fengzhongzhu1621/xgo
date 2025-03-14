package file

import (
	"os"
	"testing"
)

func Test_AtomicFileNew(t *testing.T) {
	atomFile, err := NewAtomicFile("./test.txt", os.FileMode(0755))
	if err != nil {
		t.Errorf("AtomicFileNew failed! err:%s", err.Error())
		return
	}

	defer atomFile.Close()

	_, err = atomFile.Write([]byte("test"))
	if err != nil {
		t.Errorf("fail to write data. err:%s", err.Error())
		return
	}

	// 中止写入并删除临时文件
	err = atomFile.Abort()
	if err != nil {
		t.Errorf("fail to abort file:%s", err.Error())
		return
	}
}
