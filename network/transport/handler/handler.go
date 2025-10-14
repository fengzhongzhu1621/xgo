package handler

import "context"

// Handler is the process function when server transport receive a package.
type IHandler interface {
	Handle(ctx context.Context, req []byte) (rsp []byte, err error)
}

// CloseHandler handles the logic after connection closed.
type ICloseHandler interface {
	HandleClose(ctx context.Context) error
}
