package api

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ccallazans/filedrop/internal/api/handlers"
	"github.com/ccallazans/filedrop/internal/api/middlewares"
	"github.com/ccallazans/filedrop/internal/application/usecase"
	repository "github.com/ccallazans/filedrop/internal/domain/repository/impl"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *echo.Echo {

	// Repositories
	userRepository := repository.NewUserRepository(db)
	fileRepository := repository.NewFileRepository(db)
	fileAccessRepository := repository.NewFileAccessRepository(db)

	cfg := aws.Config{
		Region: os.Getenv("AWS_REGION"),
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}
	s3client := s3.NewFromConfig(cfg)

	// Usecases
	accountUsecase := usecase.NewAccountUsecase(userRepository, fileRepository)
	uploadUsecase := usecase.NewUploadUsecase(fileRepository, fileAccessRepository, userRepository, s3client)

	// Handlers
	authHandler := handlers.NewAuthHandler(accountUsecase)
	// accountHandler := handler.NewAccountHandler(accountUsecase)
	uploadHandler := handlers.NewUploadHandler(uploadUsecase)

	// Default
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// accountGroup := e.Group("/account")
	// accountGroup.POST("")

	fileGroup := e.Group("/file")
	fileGroup.POST("/upload", middlewares.AuthenticationMiddleware(uploadHandler.UploadFile))
	fileGroup.POST("/download", middlewares.AuthenticationMiddleware(uploadHandler.AccessFile))

	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/signin", authHandler.Signin)

	return e
}
