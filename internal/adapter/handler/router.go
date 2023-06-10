package handler

import (
	"github.com/ccallazans/filedrop/internal/adapter/repository"
	"github.com/ccallazans/filedrop/internal/adapter/service"
	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *echo.Echo {

	// Repositories
	userRepository := repository.NewUserRepository(db)
	fileRepository := repository.NewFileRepository(db)
	accessFileRepository := repository.NewAccessFileRepository(db)

	// Services
	s3ClientService := service.NewS3ClientService()

	// Usecases
	accountUsecase := usecase.NewAccountUsecase(userRepository, fileRepository)
	uploadUsecase := usecase.NewUploadUsecase(userRepository, fileRepository, accessFileRepository, s3ClientService)

	// Handlers
	accountHandler := NewAccountHandler(accountUsecase)
	uploadHandler := NewUploadHandler(uploadUsecase, s3ClientService)

	// Default
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authGroup := e.Group("/auth")
	authGroup.POST("/login", accountHandler.AuthUser)
	authGroup.POST("/register", accountHandler.CreateUser)

	accessFileGroup := e.Group("/file")
	accessFileGroup.GET("/:hash", uploadHandler.AccessFile)
	accessFileGroup.POST("/upload", uploadHandler.UploadFile)

	return e
}
