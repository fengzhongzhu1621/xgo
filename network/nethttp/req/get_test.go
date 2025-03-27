package req

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/imroc/req/v3"
)

func TestSimpleGet(t *testing.T) {
	client := req.C()        // Use C() to create a client.
	resp, err := client.R(). // Use R() to create a request.
					Get("https://httpbin.org/uuid")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type UserInfo struct {
	Name string `json:"name"`
	Blog string `json:"blog"`
}

func TestAdvancedGet(t *testing.T) {
	var userInfo UserInfo
	var errMsg ErrorMessage

	// create a new client.
	client := req.C().
		SetUserAgent("my-custom-client"). // Chainable client settings.
		SetTimeout(5 * time.Second)

	// create a new request.
	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json"). // Chainable request settings.
		SetPathParam("username", "imroc").                     // Replace path variable in url.
		SetSuccessResult(&userInfo).                           // Unmarshal response body into userInfo automatically if status code is between 200 and 299.
		SetErrorResult(&errMsg).                               // Unmarshal response body into errMsg automatically if status code >= 400.
		EnableDump().                                          // Enable dump at request level, only print dump content if there is an error or some unknown situation occurs to help troubleshoot.
		Get("https://api.github.com/users/{username}")

	if err != nil { // Error handling.
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return
	}

	if resp.IsErrorState() { // Status code >= 400.
		fmt.Println(errMsg.Message) // Record error message returned.
		return
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		fmt.Printf("%s (%s)\n", userInfo.Name, userInfo.Blog)
		return
	}

	// Unknown status code.
	log.Println("unknown status", resp.Status)
	log.Println("raw content:")
	log.Println(resp.Dump()) // Record raw content when server returned unknown status code.
}

type ErrorMessage2 struct {
	Message string `json:"message"`
}

func (msg *ErrorMessage2) Error() string {
	return fmt.Sprintf("API Error: %s", msg.Message)
}

func TestOnAfterResponse(t *testing.T) {
	var userInfo UserInfo

	var client = req.C().
		SetUserAgent("my-custom-client"). // Chainable client settings.
		SetTimeout(5 * time.Second).
		EnableDumpEachRequest().
		SetCommonErrorResult(&ErrorMessage2{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
				return nil
			}
			if errMsg, ok := resp.ErrorResult().(*ErrorMessage2); ok {
				resp.Err = errMsg // Convert api error into go error
				return nil
			}
			if !resp.IsSuccessState() {
				// Neither a success response nor a error response, record details to help troubleshooting
				resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
			}
			return nil
		})

	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json"). // Chainable request settings
		SetPathParam("username", "imroc").
		SetSuccessResult(&userInfo). // Unmarshal response body into userInfo automatically if status code is between 200 and 299.
		Get("https://api.github.com/users/{username}")

	if err != nil { // Error handling.
		log.Println("error:", err)
		return
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		fmt.Printf("%s (%s)\n", userInfo.Name, userInfo.Blog)
	}
}
