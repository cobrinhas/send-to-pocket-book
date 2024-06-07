package main

import (
	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/http"
	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	as "github.com/palavrapasse/aspirador/pkg"
)

func main() {

	logging.Aspirador = as.WithClients(logging.CreateAspiradorClients(http.ServerAddress()))

	logging.Aspirador.Trace("Starting proxy-server Service")

	e := echo.New()

	defer e.Close()

	http.RegisterMiddlewares(e)
	http.RegisterHandlers(e)

	e.Logger.Fatal(http.Start(e))
}
