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
	authService service.AuthService
	fileService service.FileService
}

func NewApi(db *gorm.DB) (*api, error) {
	// Repository
	userStore := repository.NewPostgresUserStore(db)
	fileStore := repository.NewPostgresFileStore(db)

	// S3Client
	cfg, err := config.NewAWSConfig()
	if err != nil {
		return nil, err
	}

	s3Client := config.NewS3Client(cfg)

	// Services
	authService := service.NewAuthService(userStore)
	fileService := service.NewFileService(fileStore, userStore, s3Client)

	return &api{
		authService: *authService,
		fileService: *fileService,
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

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/login", a.Login)
	auth.POST("/register", a.Register)

	user := v1.Group("/users")
	user.GET("/:id", middlewares.Auth(a.FindAccountByID, []int{int(domain.ADMIN), int(domain.USER)}))

	file := v1.Group("/files")
	file.POST("/upload", middlewares.Auth(a.UploadFile, []int{}))
	file.GET("/", a.DownloadFile)

	return e
}
