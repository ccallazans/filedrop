package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/ccallazans/filedrop/internal/application/usecase"
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

	type Request struct {
		Secret string                `form:"secret" validate:"omitempty"`
		File   *multipart.FileHeader `form:"file" validate:"required"`
	}

	var request Request
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	request.File, err = c.FormFile("file")
	if err != nil {
		return err
	}

	err = validator.New().Struct(request)
	if err != nil {
		return err
	}

	hash, err := h.uploadUsecase.UploadFile(c.Request().Context(), request.Secret, request.File)
	if err != nil {
		return err
	}

	type Response struct {
		Hash string `json:"hash"`
	}

	return c.JSON(http.StatusCreated, Response{hash})
}

// func (h *UploadHandler) AccessFile(c echo.Context) error {

// 	hash := c.Param("hash")

// 	type AccessFileRequest struct {
// 		AccessCode string `json:"access_code" validate:"required"`
// 	}

// 	var accessFileRequest AccessFileRequest
// 	err := c.Bind(&accessFileRequest)
// 	if err != nil {
// 		return ErrorHandler(err, http.StatusBadRequest, c)
// 	}

// 	// err = validator.New().Struct(accessFileRequest)
// 	// if err != nil {
// 	// 	return ErrorHandler(err, http.StatusBadRequest, c)
// 	// }

// 	buffer, err := h.uploadUsecase.AccessFile(hash, usecase.AccessFileArgs{
// 		AccessCode: accessFileRequest.AccessCode,
// 	})
// 	if err != nil {
// 		return ErrorHandler(err, http.StatusBadRequest, c)
// 	}

// 	c.Response().Header().Set("Content-Type", "application/octet-stream")
// 	c.Response().Header().Set("Content-Disposition", "attachment; filename=downloaded_file.txt")

// 	_, err = c.Response().Write(buffer.Bytes())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
