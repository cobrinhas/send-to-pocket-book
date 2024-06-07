package http

import (
	"fmt"
	"net/http"
    "math/rand"

	"github.com/labstack/echo/v4"
	"github.com/cobrinhas/send-to-pocket-book/proxy-server/internal/logging"
)

func RegisterHandlers(e *echo.Echo) {

	e.GET(randomNumberRoute, GetRandomNumber)
	e.POST(randomNumberRoute, PrintRandomNumber)
	e.GET(internalServerErrorRoute, ReturnInternalServerError)

	echo.NotFoundHandler = useNotFoundHandler()
}

func GetRandomNumber(ectx echo.Context) error {

	logging.Aspirador.Trace("Getting random value")

	generatedNumber := rand.Intn(100)

	logging.Aspirador.Trace(fmt.Sprintf("Number generated %d", generatedNumber))

	return Ok(ectx, generatedNumber)
}

func PrintRandomNumber(ectx echo.Context) error {

	logging.Aspirador.Trace("Printing random value")

	receivedNumber := ectx.QueryParam(valueQueryParam)

	logging.Aspirador.Trace(fmt.Sprintf("Number received %s", receivedNumber))

	return NoContent(ectx)
}

func ReturnInternalServerError(ectx echo.Context) error {
	logging.Aspirador.Trace("Returning Internal Server Error")
	return InternalServerError(ectx)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
