package http

import (
	"fmt"
	"net/http"
	"ratelimit/pkg/echorouter"
	"ratelimit/pkg/ratelimit"

	"github.com/labstack/echo/v4"
)

// Handler http handler
type Handler struct {
	e           *echo.Echo
	ratelimiter ratelimit.Limiter
}

// NewHandler create Handler instance
func NewHandler(e *echo.Echo, ratelimiter ratelimit.Limiter) *Handler {
	h := &Handler{
		e:           e,
		ratelimiter: ratelimiter,
	}
	h.SetRoutes()
	return h
}

func (h *Handler) rateLimitEndpoint(c echo.Context) error {
	count := c.Request().Header.Get(echorouter.HeaderRateLimitRequestCount)
	realIP := c.RealIP()
	return c.String(http.StatusOK, fmt.Sprintf("IP: %s, request: %s", realIP, count))
}
