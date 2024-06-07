package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NoContent(ectx echo.Context) error {
	return ectx.NoContent(http.StatusNoContent)
}

func Ok(ectx echo.Context, i interface{}) error {
	return ectx.JSON(http.StatusOK, i)
}

func Accepted(ectx echo.Context) error {
	return ectx.NoContent(http.StatusAccepted)
}

func InternalServerError(ectx echo.Context) error {
	return ectx.NoContent(http.StatusInternalServerError)
}

func TooManyRequests(ectx echo.Context) error {
	return ectx.NoContent(http.StatusTooManyRequests)
}

func BadRequest(ectx echo.Context) error {
	return ectx.NoContent(http.StatusBadRequest)
}

func Forbidden(ectx echo.Context) error {
	return ectx.NoContent(http.StatusForbidden)
}
