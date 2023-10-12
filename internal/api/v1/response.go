package api

import (
	"mime/multipart"
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	Email     string    `json:"email"`
	Role      uint      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

//

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

//

type UploadFileRequest struct {
	Password string                `form:"password" validate:"omitempty"`
	File     *multipart.FileHeader `form:"file" validate:"required"`
}

type UploadFileResponse struct {
	Hash string `json:"hash"`
}