package resty

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestSetTimeout(t *testing.T) {
	client := resty.New()

	// 设置请求超时时间为5秒
	client.SetTimeout(5 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get("https://api.example.com/data")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode())
	fmt.Println("Response Body:", resp.String())
}
