package rest

import "net/http"

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
