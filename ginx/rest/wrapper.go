package rest

import "net/http"

// clientWrapper is a wrapper for restful client
type clientWrapper struct {
	client   ClientInterface
	wrappers RequestWrapperChain
}

// NewClientWrapper new restful client wrapper
func NewClientWrapper(client ClientInterface, wrappers ...RequestWrapper) ClientInterface {
	return &clientWrapper{
		client:   client,
		wrappers: wrappers,
	}
}

// Verb generate restful request by http method
func (r *clientWrapper) Verb(verb VerbType) *Request {
	return ProcessRequestWrapperChain(r.client.Verb(verb), r.wrappers)
}

// Post generate post restful request
func (r *clientWrapper) Post() *Request {
	return ProcessRequestWrapperChain(r.client.Post(), r.wrappers)
}

// Put generate put restful request
func (r *clientWrapper) Put() *Request {
	return ProcessRequestWrapperChain(r.client.Put(), r.wrappers)
}

// Get generate get restful request
func (r *clientWrapper) Get() *Request {
	return ProcessRequestWrapperChain(r.client.Get(), r.wrappers)
}

// Delete generate delete restful request
func (r *clientWrapper) Delete() *Request {
	return ProcessRequestWrapperChain(r.client.Delete(), r.wrappers)
}

// Patch generate patch restful request
func (r *clientWrapper) Patch() *Request {
	return ProcessRequestWrapperChain(r.client.Patch(), r.wrappers)
}

// RequestWrapper is the restful request wrapper
type RequestWrapper func(*Request) *Request

// RequestWrapperChain is the restful request wrapper chain
type RequestWrapperChain []RequestWrapper

// ProcessRequestWrapperChain process restful request wrapper chain
func ProcessRequestWrapperChain(req *Request, wrappers RequestWrapperChain) *Request {
	for _, wrapper := range wrappers {
		req = wrapper(req)
	}
	return req
}

// BaseUrlWrapper returns a restful request wrapper that changes request's base url
func BaseUrlWrapper(baseUrl string) RequestWrapper {
	return func(request *Request) *Request {
		request.baseURL = baseUrl
		return request
	}
}

// HeaderWrapper returns a restful request wrapper that changes request's header
func HeaderWrapper(handler func(http.Header) http.Header) RequestWrapper {
	return func(request *Request) *Request {
		if request.headers == nil {
			request.headers = make(http.Header)
		}
		request.headers = handler(request.headers)
		return request
	}
}
