package api

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ccallazans/filedrop/internal/api/middlewares"
	"github.com/ccallazans/filedrop/internal/application/usecase"
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

type api struct {
	logger *zap.Logger

	authUsecase usecase.AuthUsecase
	fileUsecase usecase.FileUsecase
	userUsecase usecase.UserUsecase
}

func NewApi(logger *zap.Logger, db *gorm.DB) *api {
	userStore := repository.NewPostgresUserStore(db)
	fileStore := repository.NewPostgresFileStore(db)
	fileAccessStore := repository.NewPostgresFileAccessStore(db)

	// S3Client
	cfg := aws.Config{Region: os.Getenv("AWS_REGION"), Credentials: credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")}
	s3client := s3.NewFromConfig(cfg)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userStore)
	userUsecase := usecase.NewUserUsecase(userStore, fileStore)
	fileUsecase := usecase.NewFileUsecase(fileStore, fileAccessStore, userStore, s3client)

	return &api{
		authUsecase: *authUsecase,
		fileUsecase: *fileUsecase,
		userUsecase: *userUsecase,
	}
}

func (a *api) Routes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = APIErrorHandler

	authGroup := e.Group("/auth")
	authGroup.POST("/", a.Signin)

	userGroup := e.Group("/users")
	userGroup.POST("/", a.CreateUser)
	userGroup.GET("/", middlewares.AuthenticationMiddleware(a.GetAllUsers, []int{domain.ADMIN}))
	userGroup.GET("/:id", middlewares.AuthenticationMiddleware(a.GetUserByID, []int{domain.ADMIN}))
	userGroup.DELETE("/:id", middlewares.AuthenticationMiddleware(a.DeleteUserByID, []int{domain.ADMIN}))

	fileGroup := e.Group("/file")
	fileGroup.POST("/upload", middlewares.AuthenticationMiddleware(a.UploadFile, []int{}))
	fileGroup.POST("/download", middlewares.AuthenticationMiddleware(a.AccessFile, []int{}))

	return e
}
