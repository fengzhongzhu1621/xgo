package req

import (
	"fmt"
	"log"
	"testing"

	"github.com/imroc/req/v3"
)

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Result struct {
	Data string `json:"data"`
}

func TestSimplePost(t *testing.T) {
	client := req.C().DevMode()
	var result Result

	resp, err := client.R().
		SetBody(&Repo{Name: "req", Url: "https://github.com/imroc/req"}).
		SetSuccessResult(&result).
		Post("https://httpbin.org/post")
	if err != nil {
		log.Fatal(err)
	}

	if !resp.IsSuccessState() {
		fmt.Println("bad response status:", resp.Status)
		return
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("data:", result.Data)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
}

type APIResponse struct {
	Origin string `json:"origin"`
	Url    string `json:"url"`
}

func TestDoStylePost(t *testing.T) {
	var resp APIResponse
	c := req.C().SetBaseURL("https://httpbin.org/post")
	err := c.Post().
		SetBody("hello").
		Do().
		Into(&resp)
	if err != nil {
		panic(err)
	}
	fmt.Println("My IP is", resp.Origin)
}
