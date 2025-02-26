package retry_go

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/avast/retry-go/v4"
)

func TestDo(t *testing.T) {
	url := "http://example.com"
	var body []byte

	err := retry.Do(
		func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		// 处理错误
	}

	fmt.Println(string(body))
}

func TestDoWithData(t *testing.T) {
	url := "http://example.com"

	body, err := retry.DoWithData(
		func() ([]byte, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return body, nil
		},
	)

	if err != nil {
		// 处理错误
	}

	fmt.Println(string(body))
}
