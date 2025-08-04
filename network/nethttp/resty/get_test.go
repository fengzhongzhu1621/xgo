package resty

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGet(t *testing.T) {
	client := resty.New()

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

func TestSetQueryParams(t *testing.T) {
	client := resty.New()

	// 设置查询参数
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		Get("https://api.example.com/data")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode())
	fmt.Println("Response Body:", resp.String())
}
