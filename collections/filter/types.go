package filter

import "context"

// ClientHandleFunc defines the client side filter(interceptor) function type.
type ClientHandleFunc func(ctx context.Context, req, rsp interface{}) error

// ClientFilter is the client side filter(interceptor) type. They are chained to process request.
type ClientFilter func(ctx context.Context, req, rsp interface{}, next ClientHandleFunc) error

// ServerHandleFunc defines the server side filter(interceptor) function type.
type ServerHandleFunc func(ctx context.Context, req interface{}) (rsp interface{}, err error)

// ServerFilter is the server side filter(interceptor) type. They are chained to process request.
type ServerFilter func(ctx context.Context, req interface{}, next ServerHandleFunc) (rsp interface{}, err error)
