package auth

import (
	"context"
	"net/http"
)

type ISecurityProvider interface {
	Intercept(ctx context.Context, req *http.Request) error
}
