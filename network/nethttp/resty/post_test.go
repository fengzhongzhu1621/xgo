package resty

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestPost(t *testing.T) {
	client := resty.New()

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		ForceContentType("application/json").
		SetBody(data).
		Post("https://api.example.com/data")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode())
	fmt.Println("Response Body:", resp.String())
}
