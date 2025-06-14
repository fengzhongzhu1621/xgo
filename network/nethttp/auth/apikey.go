package auth

import (
	"context"
	"net/http"
)

// 这是一个预定义的错误，表示在使用 SecurityProviderApiKey 时，指定的 in 参数值无效。有效的 in 参数应为 "cookie"、"header" 或 "query"。
const (
	// ErrSecurityProviderApiKeyInvalidIn indicates a usage of an invalid In.
	// Should be cookie, header or query
	ErrSecurityProviderApiKeyInvalidIn = SecurityProviderError("invalid 'in' specified for apiKey")
)

// SecurityProviderError defines error values of a security provider.
type SecurityProviderError string

// Error implements the error interface.
func (e SecurityProviderError) Error() string {
	return string(e)
}

// NewSecurityProviderApiKey will attach a generic apiKey for a given name
// either to a cookie, header or as a query parameter.
func NewSecurityProviderApiKey(in, name, apiKey string) (*SecurityProviderApiKey, error) {
	interceptors := map[string]func(ctx context.Context, req *http.Request) error{
		"cookie": func(ctx context.Context, req *http.Request) error {
			req.AddCookie(&http.Cookie{Name: name, Value: apiKey})
			return nil
		},
		"header": func(ctx context.Context, req *http.Request) error {
			req.Header.Add(name, apiKey)
			return nil
		},
		"query": func(ctx context.Context, req *http.Request) error {
			query := req.URL.Query()
			query.Add(name, apiKey)
			req.URL.RawQuery = query.Encode()
			return nil
		},
	}

	interceptor, ok := interceptors[in]
	if !ok {
		return nil, ErrSecurityProviderApiKeyInvalidIn
	}

	return &SecurityProviderApiKey{
		interceptor: interceptor,
	}, nil
}

// SecurityProviderApiKey will attach an apiKey either to a
// cookie, header or query.
type SecurityProviderApiKey struct {
	interceptor func(ctx context.Context, req *http.Request) error
}

// Intercept will attach a cookie, header or query param for the configured
// name and apiKey.
func (s *SecurityProviderApiKey) Intercept(ctx context.Context, req *http.Request) error {
	return s.interceptor(ctx, req)
}
