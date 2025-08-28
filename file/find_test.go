package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocateFile(t *testing.T) {
	pwd, _ := os.Getwd()
	tests := []struct {
		name     string
		filename string
		dirs     []string
		expect   string
	}{
		{name: "1", filename: "file.go", dirs: []string{pwd}, expect: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := LocateFile(tt.filename, tt.dirs)
			fmt.Println(actual)
		})
	}
}

func TestWhich(t *testing.T) {
	filepath, _ := Which("sh")
	assert.Equal(t, filepath, "/bin/sh")

	filepath, err := Which("xxx")
	assert.Equal(t, filepath, "")
	assert.ErrorContains(t, err, "no such file or directory")
}
