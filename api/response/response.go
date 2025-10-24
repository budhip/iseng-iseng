package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DON'T CHANGE THIS STRUCT
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// DON'T CHANGE THIS RESPONSE
// Success sends a success response with the given status code and data.
func Success(c echo.Context, v any) error {
	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    v,
	})
}

// DON'T CHANGE THIS RESPONSE
// Created sends a created response with the given status code and data.
func Created(c echo.Context, v any) error {
	return c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    v,
	})
}

// DON'T CHANGE THIS RESPONSE
// Error sends an error response with the given status code and message.
func Error(c echo.Context, status int, message string) error {
	return c.JSON(status, Response{
		Status:  status,
		Message: message,
	})
}
