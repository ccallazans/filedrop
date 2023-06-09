package handler

import (
	"github.com/ccallazans/filedrop/internal/adapter/repository"
	"github.com/ccallazans/filedrop/internal/adapter/service"
	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *echo.Echo {

	// Repositories
	userRepository := repository.NewUserRepository(db)
	fileRepository := repository.NewFileRepository(db)

	// Services
	s3ClientService := service.NewS3ClientService()

	// Usecases
	accountUsecase := usecase.NewAccountUsecase(userRepository, fileRepository)
	_ = usecase.NewUploadUsecase(userRepository, fileRepository, s3ClientService)

	// Handlers
	accountHandler := NewAccountHandler(accountUsecase)

	// Default
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authGroup := e.Group("/auth")
	authGroup.POST("/login", accountHandler.AuthUser)
	authGroup.POST("/register", accountHandler.CreateUser)

	return e
}
