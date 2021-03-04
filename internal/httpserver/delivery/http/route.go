package http

import (
	"ratelimit/pkg/echorouter"
)

// SetRoutes ...
func (h *Handler) SetRoutes() {
	basic := h.e.Group("", echorouter.MiddlewareRateLimit(h.ratelimiter))
	basic.GET("/ratelimit", h.rateLimitEndpoint)
}
