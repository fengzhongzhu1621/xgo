package nethttp

import (
	"fmt"
	"log"
	"testing"

	"github.com/duke-git/lancet/v2/netutil"
)

// TestHttpRequest HttpRequest is a struct used to abstract HTTP request entity.
//
// type HttpRequest struct {
// 	RawURL      string
// 	Method      string
// 	Headers     http.Header
// 	QueryParams url.Values
// 	FormData    url.Values
// 	File        *File
// 	Body        []byte
// }

func TestHttpRequest(t *testing.T) {
	// header := http.Header{}
	// header.Add("Content-Type", "multipart/form-data")

	// postData := url.Values{}
	// postData.Add("userId", "1")
	// postData.Add("title", "testItem")

	// request := &netutil.HttpRequest{
	// 	RawURL:   "https://jsonplaceholder.typicode.com/todos",
	// 	Method:   "POST",
	// 	Headers:  header,
	// 	FormData: postData,
	// }
}

// TestNewHttpClientWithConfig Encode url query string values.
// func EncodeUrl(urlStr string) (string, error)
//
//	type HttpClientConfig struct {
//		Timeout          time.Duration
//		SSLEnabled       bool
//		TLSConfig        *tls.Config
//		Compressed       bool
//		HandshakeTimeout time.Duration
//		ResponseTimeout  time.Duration
//		Verbose          bool
//		Proxy            *url.URL
//	}
//
//	type HttpClient struct {
//		*http.Client
//		TLS     *tls.Config
//		Request *http.Request
//		Config  HttpClientConfig
//		Context context.Context
//	}
func TestNewHttpClientWithConfig(t *testing.T) {
	// httpClientCfg := netutil.HttpClientConfig{
	// 	SSLEnabled:       true,
	// 	HandshakeTimeout: 10 * time.Second,
	// }
	// httpClient := netutil.NewHttpClientWithConfig(&httpClientCfg)
}

// TestSendRequest Use HttpClient to send HTTP request.
// func (client *HttpClient) SendRequest(request *HttpRequest) (*http.Response, error)
//
// Decode http response into target object.
// func (client *HttpClient) DecodeResponse(resp *http.Response, target any) error
func TestSendGetRequest(t *testing.T) {
	request := &netutil.HttpRequest{
		RawURL: "https://jsonplaceholder.typicode.com/todos/1",
		Method: "GET",
	}

	httpClient := netutil.NewHttpClient()
	resp, err := httpClient.SendRequest(request)
	if err != nil || resp.StatusCode != 200 {
		return
	}

	type Todo struct {
		UserId    int    `json:"userId"`
		Id        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	var todo Todo
	err = httpClient.DecodeResponse(resp, &todo)
	if err != nil {
		return
	}

	fmt.Println(todo.Id)
}

// TestParseHttpResponse Decode http response to specified interface.
// 和 httpClient.DecodeResponse 功能一样
// func ParseHttpResponse(resp *http.Response, obj any) error
func TestParseHttpResponse(t *testing.T) {
	url := "https://jsonplaceholder.typicode.com/todos/1"
	header := map[string]string{
		"Content-Type": "application/json",
	}

	resp, err := netutil.HttpGet(url, header)
	if err != nil {
		log.Fatal(err)
	}

	type Todo struct {
		Id        int    `json:"id"`
		UserId    int    `json:"userId"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	toDoResp := &Todo{}
	err = netutil.ParseHttpResponse(resp, toDoResp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(toDoResp.Id)
}

// TestStructToUrlValues Convert struct to url values, only convert the field which is exported and has `json` tag.
// func StructToUrlValues(targetStruct any) url.Values
func TestStructToUrlValues(t *testing.T) {
	type TodoQuery struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	todoQuery := TodoQuery{
		Id:   1,
		Name: "Test",
	}
	todoValues, _ := netutil.StructToUrlValues(todoQuery)

	fmt.Println(todoValues.Get("id"))   // 1
	fmt.Println(todoValues.Get("name")) //Test
}
