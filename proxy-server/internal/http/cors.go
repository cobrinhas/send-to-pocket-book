package http

import (
	"os"

	"github.com/labstack/echo/v4/middleware"
)

const (
	allowOriginsEnvKey = "allow_origins"
)

var (
	allowOrigins = os.Getenv(allowOriginsEnvKey)
)

func CORSConfig() middleware.CORSConfig {
	cc := middleware.DefaultCORSConfig

	if len(allowOrigins) > 0 {
		cc.AllowOrigins = []string{allowOrigins}
	}

	return cc
}
