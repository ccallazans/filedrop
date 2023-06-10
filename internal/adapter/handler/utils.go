package handler

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ErrorHandler(err error, status int, c echo.Context) error {
	// Create the error response message
	errorResponse := ErrorResponse{
		Message: err.Error(),
	}

	// Send the error response to the client
	return c.JSON(status, errorResponse)
}
