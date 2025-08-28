package auth

import (
	"context"
	"fmt"
	"net/http"
)

// SecurityProviderBearerToken sends a token as part of an
// Authorization: Bearer header along with a request.
// 使用Bearer Token认证
type SecurityProviderBearerToken struct {
	token string
}

// NewSecurityProviderBearerToken provides a SecurityProvider, which can solve
// the Bearer Auth challenge for api-calls.
func NewSecurityProviderBearerToken(token string) (*SecurityProviderBearerToken, error) {
	return &SecurityProviderBearerToken{
		token: token,
	}, nil
}

// Intercept will attach an Authorization header to the request
// and ensures that the bearer token is attached to the header.
func (s *SecurityProviderBearerToken) Intercept(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	return nil
}
