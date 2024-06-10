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

func BadRequest(ectx echo.Context, i interface{}) error {
	return ectx.JSON(http.StatusBadRequest, i)
}

func Forbidden(ectx echo.Context) error {
	return ectx.NoContent(http.StatusForbidden)
}

type InvalidRequestFieldResponse struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func InvalidEmailResponse() InvalidRequestFieldResponse {
	return InvalidRequestFieldResponse{
		Field:  "email",
		Reason: "invalid",
	}
}

func UnsupportedEmailDomainResponse() InvalidRequestFieldResponse {
	return InvalidRequestFieldResponse{
		Field:  "email",
		Reason: "unsupported-domain",
	}
}

func NotConnectableUrlResponse() InvalidRequestFieldResponse {
	return InvalidRequestFieldResponse{
		Field:  "url",
		Reason: "not-connectable",
	}
}

func UnsupportedDocumentResponse() InvalidRequestFieldResponse {
	return InvalidRequestFieldResponse{
		Field:  "url",
		Reason: "unsupported-document",
	}
}
