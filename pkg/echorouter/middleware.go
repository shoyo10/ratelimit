package echorouter

import (
	"fmt"
	"net/http"
	"ratelimit/pkg/errors"
	"ratelimit/pkg/ratelimit"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// MiddlewareCorsConfig ...
var MiddlewareCorsConfig = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	},
	AllowHeaders: []string{
		"*",
	},
	ExposeHeaders: []string{
		"*",
	},
})

func requestIDGenerator() string {
	return uuid.New().String()
}

// MiddlewareRequestID middleware to set request id
func MiddlewareRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = requestIDGenerator()
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

// MiddlewareLogWithRequestID set request to log field
func MiddlewareLogWithRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = res.Header().Get(echo.HeaderXRequestID)
				if rid == "" {
					rid = requestIDGenerator()
					res.Header().Set(echo.HeaderXRequestID, rid)
				}
				logger := log.With().Str("request_id", rid).Logger()
				ctx := logger.WithContext(req.Context())
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}

// MiddlewareRateLimit rate limit middleware
func MiddlewareRateLimit(limiter ratelimit.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			realIP := c.RealIP()
			resp, err := limiter.Avaliable(c.Request().Context(), realIP)
			if err != nil {
				return errors.Wrap(errors.ErrInternalServerError, err.Error())
			}
			if !resp.IsAvaliable {
				return c.String(http.StatusTooManyRequests, fmt.Sprintf("Error: [%s] too many request", realIP))
			}
			c.Request().Header.Add(HeaderRateLimitRequestCount, fmt.Sprintf("%d", resp.Count))
			return next(c)
		}
	}
}
