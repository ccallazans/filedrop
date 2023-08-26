package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uploadUsecase usecase.UploadUsecase
}

func NewUploadHandler(uploadUsecase usecase.UploadUsecase) *UploadHandler {
	return &UploadHandler{
		uploadUsecase: uploadUsecase,
	}
}

func (h *UploadHandler) UploadFile(c echo.Context) error {

	type UploadFileRequest struct {
		Secret string                `form:"secret" validate:"omitempty"`
		File   *multipart.FileHeader `form:"file" validate:"required"`
	}

	var request UploadFileRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	request.File, err = c.FormFile("file")
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	hash, err := h.uploadUsecase.UploadFile(c.Request().Context(), request.Secret, request.File)
	if err != nil {
		return ParseApiError(err)
	}

	type UploadFileResponse struct {
		Hash string `json:"hash"`
	}

	return c.JSON(http.StatusCreated, UploadFileResponse{hash})
}

func (h *UploadHandler) AccessFile(c echo.Context) error {

	type AccessFileRequest struct {
		Hash   string `json:"hash" validate:"required"`
		Secret string `json:"secret" validate:"omitempty"`
	}

	var request AccessFileRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.ErrorType{Type: utils.BadRequestErr, Message: err.Error()})
	}

	file, err := h.uploadUsecase.AccessFile(c.Request().Context(), request.Hash, request.Secret)
	if err != nil {
		return ParseApiError(err)
	}
	defer file.Body.Close()

	c.Response().Header().Set("Content-Type", *file.ContentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", file.ContentLength))
	io.Copy(c.Response().Writer, file.Body)

	return nil
}
