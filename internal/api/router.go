package api

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	handler "github.com/ccallazans/filedrop/internal/api/handlers"
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

	// S3Client
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("error creating s3 client")
	}
	s3Client := s3.New(sess, aws.NewConfig().WithRegion("sa-east-1"))

	// Usecases
	// accountUsecase := usecase.NewAccountUsecase(userRepository, fileRepository)
	uploadUsecase := usecase.NewUploadUsecase(fileRepository, fileAccessRepository, userRepository, s3Client)

	// Handlers
	// accountHandler := handler.NewUploadHandler(uploadUsecase)
	uploadHandler := handler.NewUploadHandler(uploadUsecase)

	// Default
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fileGroup := e.Group("/file")
	// fileGroup.GET("/:hash", uploadHandler.)
	fileGroup.POST("/upload", uploadHandler.UploadFile)

	return e
}
