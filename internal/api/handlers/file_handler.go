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

type FileHandler struct {
	fileUsecase usecase.FileUsecase
}

func NewFileHandler(fileUsecase usecase.FileUsecase) *FileHandler {
	return &FileHandler{
		fileUsecase: fileUsecase,
	}
}

func (h *FileHandler) UploadFile(c echo.Context) error {

	type UploadFileRequest struct {
		Secret string                `form:"secret" validate:"omitempty"`
		File   *multipart.FileHeader `form:"file" validate:"required"`
	}

	var request UploadFileRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(&utils.BadRequestError{})
	}

	request.File, err = c.FormFile("file")
	if err != nil {
		return ParseApiError(&utils.BadRequestError{})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.BadRequestError{})
	}

	hash, err := h.fileUsecase.UploadFile(c.Request().Context(), request.Secret, request.File)
	if err != nil {
		return ParseApiError(err)
	}

	type UploadFileResponse struct {
		Hash string `json:"hash"`
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/file/download/%s", hash))
	return c.JSON(http.StatusCreated, UploadFileResponse{hash})
}

func (h *FileHandler) AccessFile(c echo.Context) error {

	type AccessFileRequest struct {
		Hash   string `json:"hash" validate:"required"`
		Secret string `json:"secret" validate:"omitempty"`
	}

	var request AccessFileRequest
	err := c.Bind(&request)
	if err != nil {
		return ParseApiError(&utils.ValidationError{Message: "bad request"})
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ParseApiError(&utils.ValidationError{Message: "bad request"})
	}

	file, err := h.fileUsecase.DownloadFile(c.Request().Context(), request.Hash, request.Secret)
	if err != nil {
		return ParseApiError(err)
	}
	defer file.Body.Close()

	c.Response().Header().Set("Content-Type", *file.ContentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", file.ContentLength))
	c.Response().Header().Set("Content-Disposition", "attachment;")
	io.Copy(c.Response().Writer, file.Body)

	return nil
}
