package fileutils

import (
	"io/ioutil"
	"testing"
)

// TempMkdir makes a temporary directory
func TempMkdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "fsnotify")
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir
}

// TempMkFile makes a temporary file.
func TempMkFile(t *testing.T, dir string) string {
	f, err := ioutil.TempFile(dir, "fsnotify")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer f.Close()
	return f.Name()
}
