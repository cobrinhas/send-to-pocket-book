package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	
	"github.com/labstack/echo/v4/middleware"
	
	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(loggingMiddleware())
	
	e.Use(corsMiddleware())
	
	
	e.Use(throttlingMiddleware())
	
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

			if !te.CanAllowRequest(*req) {
				logging.Aspirador.Warning(fmt.Sprintf("denying request for %s endpoint, too many requests (429)", req.URL.Path))

				return TooManyRequests(ectx)
			}

			return next(ectx)
		}
	}
}
