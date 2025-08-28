package middlewares

import (
	"context"
	"fmt"

	"github.com/go-dew/dew"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

var (
	AdminID  = 1
	MemberID = 2
)

type userCtxKey struct{}

type CurrentUser struct {
	ID int
}

func AuthContext(ctx context.Context, u *CurrentUser) context.Context {
	return context.WithValue(ctx, userCtxKey{}, u)
}

func getCurrentUser(ctx context.Context) *CurrentUser {
	return ctx.Value(userCtxKey{}).(*CurrentUser)
}

// isAdmin checks if the current user is authorized.
func isAdmin(ctx context.Context) bool {
	return getCurrentUser(ctx).ID == AdminID
}

func AdminOnly(next dew.Middleware) dew.Middleware {
	return dew.MiddlewareFunc(func(ctx dew.Context) error {
		if !isAdmin(ctx.Context()) {
			// Return an unauthorized error.
			return ErrUnauthorized
		}
		// Continue to the next middleware or handler.
		return next.Handle(ctx)
	})
}
