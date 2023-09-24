package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (a *api) UploadFile(c echo.Context) error {
	type UploadFileRequest struct {
		Password string                `form:"password" validate:"omitempty"`
		File     *multipart.FileHeader `form:"file" validate:"required"`
	}

	var request UploadFileRequest
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

	hash, err := a.fileService.Upload(c.Request().Context(), request.Password, request.File)
	if err != nil {
		return err
	}

	type UploadFileResponse struct {
		Hash string `json:"hash"`
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/file/download/%s", hash))
	return c.JSON(http.StatusCreated, UploadFileResponse{hash})
}

// TODO: use queryParam to get file
func (a *api) DownloadFile(c echo.Context) error {
	hash := c.QueryParam("hash")
	key := c.QueryParam("key")

	file, filename, err := a.fileService.DownloadFile(c.Request().Context(), hash, key)
	if err != nil {
		return err
	}
	defer file.Body.Close()

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Response().Header().Set("Content-Type", *file.ContentType)
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", file.ContentLength))
	io.Copy(c.Response().Writer, file.Body)

	return nil
}
