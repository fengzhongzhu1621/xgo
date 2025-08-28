package process

import (
	"os"
	"testing"
)

func TestSaveAndReadPid(t *testing.T) {
	curPid := os.Getpid()

	if err := SavePid(); err != nil {
		t.Errorf("fail to save pid. err:%s", err.Error())
	}

	pid, err := ReadPid()
	if err != nil {
		t.Errorf("fail to read pid. err:%s", err.Error())
	}

	if pid != curPid {
		t.Errorf("pid: %d, but want: %d", pid, curPid)
	}
}
