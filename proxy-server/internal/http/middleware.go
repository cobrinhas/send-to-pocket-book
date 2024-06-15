package http

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
)

const (
	securityTokenEnvKey = "security_token"
)

var (
	securityToken = os.Getenv(securityTokenEnvKey)
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(loggingMiddleware())
	e.Use(corsMiddleware())
	e.Use(throttlingMiddleware())
	e.Use(validateSecurityTokenMiddleware())

}

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {

			req := ectx.Request()

			logging.Aspirador.Info(fmt.Sprintf("Host: %s | Method: %s | Path: %s | Client IP: %s", req.Host, req.Method, req.URL.RequestURI(), ectx.RealIP()))

			return next(ectx)
		}
	}
}

func corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(CORSConfig())
}

func throttlingMiddleware() echo.MiddlewareFunc {
	te := NewThrottlingEngine()

	te.StartThrottlingEngineCleanUp()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			req := ectx.Request()

			if !te.CanAllowRequest(*req, ectx) {
				logging.Aspirador.Warning(fmt.Sprintf("denying request for %s endpoint, too many requests (429)", req.URL.Path))

				return TooManyRequests(ectx)
			}

			return next(ectx)
		}
	}
}

func validateSecurityTokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			req := ectx.Request()
			reqSecurityToken := req.Header.Get(securityTokenHeader)

			if reqSecurityToken != securityToken {
				logging.Aspirador.Error("security tokens do not match!")

				return Forbidden(ectx)
			}

			return next(ectx)
		}
	}
}
