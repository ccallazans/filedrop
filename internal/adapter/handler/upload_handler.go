package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/ccallazans/filedrop/internal/application/service"
	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uploadUsecase usecase.UploadUsecase
	s3Client      service.IS3Client
}

func NewUploadHandler(uploadUsecase usecase.UploadUsecase, s3Client service.IS3Client) *UploadHandler {
	return &UploadHandler{
		uploadUsecase: uploadUsecase,
		s3Client:      s3Client,
	}
}

func (h *UploadHandler) UploadFile(c echo.Context) error {

	type UploadFileRequest struct {
		Lock       bool                  `form:"lock"`
		AccessCode string                `form:"access_code"`
		File       *multipart.FileHeader `form:"file"`
	}

	var request UploadFileRequest
	err := c.Bind(&request)
	if err != nil {
		return ErrorHandler(err, http.StatusBadRequest, c)
	}

	request.File, err = c.FormFile("file")
	if err != nil {
		return ErrorHandler(err, http.StatusBadRequest, c)
	}

	err = validator.New().Struct(request)
	if err != nil {
		return ErrorHandler(err, http.StatusBadRequest, c)
	}

	err = h.uploadUsecase.UploadFile(&usecase.UploadFileArgs{
		Lock:       request.Lock,
		AccessCode: request.AccessCode,
		File:       request.File,
	})
	if err != nil {
		return ErrorHandler(err, http.StatusBadRequest, c)
	}

	return nil
}

func (h *UploadHandler) AccessFile(c echo.Context) error {

	hash := c.Param("hash")

	type AccessFileRequest struct {
		AccessCode string `json:"access_code" validate:"required"`
	}

	var accessFileRequest AccessFileRequest
	err := c.Bind(&accessFileRequest)
	if err != nil {
		return ErrorHandler(err, http.StatusBadRequest, c)
	}

	// err = validator.New().Struct(accessFileRequest)
	// if err != nil {
	// 	return ErrorHandler(err, http.StatusBadRequest, c)
	// }

	buffer, err := h.uploadUsecase.AccessFile(hash, usecase.AccessFileArgs{
		AccessCode: accessFileRequest.AccessCode,
	})
	if err != nil {
		return ErrorHandler(err, http.StatusBadRequest, c)
	}

	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=downloaded_file.txt")

	_, err = c.Response().Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
