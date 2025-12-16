package app

import (
	"context"
	"os"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/frontend"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/domain/mood"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type App struct {
	authService *auth.Service
	echo        *echo.Echo
	config      Config
}

func New() *App {
	if strconv.IntSize != 64 {
		logging.Error(context.Background(), "unsupported architecture: only 64-bit systems are supported")
		os.Exit(1)
	}

	_ = godotenv.Load()

	config, err := GetConfig()
	if err != nil {
		logging.Error(context.Background(), "failed to get configuration", "error", err)
		os.Exit(1)
	}

	logging.Setup(config.Logging)

	database, err := db.Connect(config.DB)
	if err != nil {
		logging.Error(context.Background(), "failed to connect to database", "error", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(database,
		mood.Entry{},
		auth.RefreshToken{},
	)
	if err != nil {
		logging.Error(context.Background(), "database migration failed", "error", err)
		os.Exit(1)
	}

	e := server.NewEchoServer(config.Server)

	frontend.RegisterRoutes(e, config.Frontend)

	authRepository := auth.NewRepository(database)
	authService := auth.NewService(authRepository, config.Auth, config.StubUserConfig)
	authHandler := auth.NewHandler(authService, config.Auth, config.StubUserConfig)

	jwtMiddleware := auth.JWTMiddleware(config.Auth, authService)

	auth.RegisterRoutes(e, authHandler, jwtMiddleware)

	moodRepo := mood.NewRepository(database)
	moodService := mood.NewService(moodRepo)
	moodHandler := mood.NewHandler(moodService)

	mood.RegisterRoutes(e, moodHandler, jwtMiddleware)

	return &App{
		authService: authService,
		echo:        e,
		config:      config,
	}
}

func (a *App) Start() {
	if err := a.authService.DeleteExpiredRefreshTokens(context.Background()); err != nil {
		logging.Error(context.Background(), "failed to delete expired refresh tokens", "error", err)
		os.Exit(1)
	}

	if err := server.StartServer(a.echo, a.config.Server); err != nil {
		logging.Error(context.Background(), "general server failure", "error", err)
		os.Exit(1)
	}
}
