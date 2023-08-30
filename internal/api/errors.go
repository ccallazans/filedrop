package api

import (
	"net/http"

	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func APIErrorHandler(err error, c echo.Context) {
	switch err.(type) {
	case *utils.AuthenticationError:
		response(c, err, http.StatusUnauthorized, "Unauthorized", err.Error())

	case *utils.AuthorizationError:
		response(c, err, http.StatusForbidden, "Forbidden", err.Error())

	case *utils.BadRequestError:
		message := err.Error()
		if len(message) == 0 {
			message = "bad request"
		}
		response(c, err, http.StatusBadRequest, "Bad Request", message)

	case *utils.ValidationError:
		response(c, err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error())

	case *utils.ConflictError:
		response(c, err, http.StatusConflict, "Conflict", err.Error())

	case *utils.NotFoundError:
		response(c, err, http.StatusNotFound, "Not Found", err.Error())

	default:
		response(c, err, http.StatusInternalServerError, "Internal server error", err.Error())
	}
}

func response(c echo.Context, err error, statusCode int, title string, detail string) {
	c.JSON(statusCode, errorResponse{
		Status: statusCode,
		Title:  title,
		Detail: err.Error(),
	})
}
