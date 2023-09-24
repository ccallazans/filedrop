package api

import (
	"github.com/ccallazans/filedrop/internal/api/v1/middlewares"
	"github.com/ccallazans/filedrop/internal/application/service"

	"github.com/ccallazans/filedrop/internal/config"
	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/ccallazans/filedrop/internal/domain/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

type api struct {
	accountService service.AccountService
	fileService    service.FileService
}

func NewApi(db *gorm.DB) (*api, error) {
	// Repository
	userStore := repository.NewPostgresUserStore(db)
	fileStore := repository.NewPostgresFileStore(db)

	// S3Client
	awsConfig, err := config.NewAWSConfig()
	if err != nil {
		return nil, err
	}
	s3Client := config.NewS3Client(awsConfig)

	// Services
	accountService := service.NewAccountService(userStore)
	fileService := service.NewFileService(fileStore, userStore, s3Client)

	return &api{
		accountService: *accountService,
		fileService:    *fileService,
	}, nil
}

func (a *api) Routes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.HTTPErrorHandler = APIErrorHandler

	account := e.Group("/accounts")
	account.POST("/login", a.Login)
	account.POST("/register", a.Register)
	account.GET("/:id", middlewares.Auth(a.FindAccountByID, []int{int(domain.ADMIN), int(domain.USER)}))

	fileGroup := e.Group("/files")
	fileGroup.POST("/upload", middlewares.Auth(a.UploadFile, []int{}))
	// fileGroup.GET("/download", middlewares.Auth(a.DownloadFile, []int{}))
	fileGroup.GET("/download", a.DownloadFile)

	return e
}
