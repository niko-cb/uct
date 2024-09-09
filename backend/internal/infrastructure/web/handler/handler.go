package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/niko-cb/uct/internal/infrastructure/web/config/actx"
)

func withContext(c echo.Context, fn func(ctx context.Context) error) error {

	// Add the echo context to the new context
	ctx := actx.Context(context.Background(), actx.EchoContext, c)

	// Replace the request's context with the new context
	req := c.Request().WithContext(ctx)
	c.SetRequest(req)

	// Call the provided function with the new context
	return fn(ctx)
}
