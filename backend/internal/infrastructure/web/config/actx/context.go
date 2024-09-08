package actx

import (
	"context"
)

type contextKey string

const (
	EchoContext contextKey = "EchoContext"
)

// Get retrieves a value from the context
func Get(ctx context.Context, key contextKey) interface{} {
	return ctx.Value(key)
}

// Context sets a value in the context
func Context(ctx context.Context, key contextKey, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}
