package auth

import (
	"context"
	"net/http"
)

// SecurityProviderBasicAuth sends a base64-encoded combination of
// username, password along with a request.
type SecurityProviderBasicAuth struct {
	username string
	password string
}

// NewSecurityProviderBasicAuth provides a SecurityProvider, which can solve
// the BasicAuth challenge for api-calls.
func NewSecurityProviderBasicAuth(username, password string) (*SecurityProviderBasicAuth, error) {
	return &SecurityProviderBasicAuth{
		username: username,
		password: password,
	}, nil
}

// Intercept will attach an Authorization header to the request and ensures that
// the username, password are base64 encoded and attached to the header.
func (s *SecurityProviderBasicAuth) Intercept(ctx context.Context, req *http.Request) error {
	req.SetBasicAuth(s.username, s.password)
	return nil
}
