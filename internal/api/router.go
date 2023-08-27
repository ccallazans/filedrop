package api

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ccallazans/filedrop/internal/api/handlers"
	"github.com/ccallazans/filedrop/internal/api/middlewares"
	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *echo.Echo {

	// Repositories
	userStore := repository.NewPostgresUserStore(db)
	fileStore := repository.NewPostgresFileStore(db)
	fileAccessStore := repository.NewPostgresFileAccessStore(db)

	// S3Client
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
	authUsecase := usecase.NewAuthUsecase(userStore)
	userUsecase := usecase.NewUserUsecase(userStore, fileStore)
	fileUsecase := usecase.NewFileUsecase(fileStore, fileAccessStore, userStore, s3client)

	// Handlers
	authHandler := handlers.NewAuthHandler(*authUsecase)
	userHandler := handlers.NewUserHandler(*userUsecase)
	uploadHandler := handlers.NewFileHandler(*fileUsecase)

	// Default
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authGroup := e.Group("/auth")
	authGroup.POST("/", authHandler.Signin)

	userGroup := e.Group("/users")
	userGroup.POST("/", userHandler.CreateUser)
	userGroup.GET("/", middlewares.AuthenticationMiddleware(userHandler.GetAllUsers, []int{domain.ADMIN}))
	userGroup.GET("/:id", middlewares.AuthenticationMiddleware(userHandler.GetUserByID, []int{domain.ADMIN}))
	userGroup.DELETE("/:id", middlewares.AuthenticationMiddleware(userHandler.DeleteUserByID, []int{domain.ADMIN}))

	fileGroup := e.Group("/file")
	fileGroup.POST("/upload", middlewares.AuthenticationMiddleware(uploadHandler.UploadFile, []int{}))
	fileGroup.POST("/download", middlewares.AuthenticationMiddleware(uploadHandler.AccessFile, []int{}))

	return e
}
