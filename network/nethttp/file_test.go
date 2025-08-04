package nethttp

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/netutil"
)

// TestDownloadFile Download the file exist in url to a local file.
// func DownloadFile(filepath string, url string) error
func TestDownloadFile(t *testing.T) {
	err := netutil.DownloadFile(
		"./lancet_logo.jpg",
		"https://picx.zhimg.com/v2-fc82a4199749de9cfb71e32e54f489d3_720w.jpg?source=172ae18b",
	)

	fmt.Println(err)
}

// TestUploadFile Upload the file to a server.
// func UploadFile(filepath string, server string) (bool, error)
func TestUploadFile(t *testing.T) {
	ok, err := netutil.UploadFile("./a.jpg", "http://www.xxx.com/bucket/test")

	fmt.Println(ok)
	fmt.Println(err)
}
