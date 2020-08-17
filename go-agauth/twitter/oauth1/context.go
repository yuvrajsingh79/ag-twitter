package oauth1

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey struct{}

// unexported key type prevents collisions
type key int

const (
	errorKey        key = iota
	requestTokenKey key = iota
	requestSecretKey
	accessTokenKey
	accessSecretKey
)

// HTTPClient is the context key to associate an *http.Client value with
// a context.
var HTTPClient contextKey

// NoContext is the default context to use in most cases.
var NoContext = context.TODO()

// contextTransport gets the Transport from the context client or nil.
func contextTransport(ctx context.Context) http.RoundTripper {
	if client, ok := ctx.Value(HTTPClient).(*http.Client); ok {
		return client.Transport
	}
	return nil
}

// WithError returns a copy of ctx that stores the given error value.
func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}

// ErrorFromContext returns the error value from the ctx or an error that the
// context was missing an error value.
func ErrorFromContext(ctx context.Context) error {
	err, ok := ctx.Value(errorKey).(error)
	if !ok {
		return fmt.Errorf("Context missing error value")
	}
	return err
}

// WithRequestToken returns a copy of ctx that stores the request token and
// secret values.
func WithRequestToken(ctx context.Context, requestToken, requestSecret string) context.Context {
	ctx = context.WithValue(ctx, requestTokenKey, requestToken)
	ctx = context.WithValue(ctx, requestSecretKey, requestSecret)
	return ctx
}

// RequestTokenFromContext returns the request token and secret from the ctx.
func RequestTokenFromContext(ctx context.Context) (string, string, error) {
	requestToken, okT := ctx.Value(requestTokenKey).(string)
	requestSecret, okS := ctx.Value(requestSecretKey).(string)
	if !okT || !okS {
		return "", "", fmt.Errorf("oauth1: Context missing request token or secret")
	}
	return requestToken, requestSecret, nil
}

// WithAccessToken returns a copy of ctx that stores the access token and
// secret values.
func WithAccessToken(ctx context.Context, accessToken, accessSecret string) context.Context {
	ctx = context.WithValue(ctx, accessTokenKey, accessToken)
	ctx = context.WithValue(ctx, accessSecretKey, accessSecret)
	return ctx
}

// AccessTokenFromContext returns the access token and secret from the ctx.
func AccessTokenFromContext(ctx context.Context) (string, string, error) {
	accessToken, okT := ctx.Value(accessTokenKey).(string)
	accessSecret, okS := ctx.Value(accessSecretKey).(string)
	if !okT || !okS {
		return "", "", fmt.Errorf("oauth1: Context missing access token or secret")
	}
	return accessToken, accessSecret, nil
}
