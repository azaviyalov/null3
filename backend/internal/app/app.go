package app

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/frontend"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/admin"
	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type App struct {
	sessionService *session.Service
	echo           *echo.Echo
	config         Config
}

func New() *App {
	if strconv.IntSize != 64 {
		slog.Error("unsupported architecture: only 64-bit systems are supported")
		os.Exit(1)
	}

	_ = godotenv.Load()
	logging.Setup(logging.GetConfig())

	config, err := GetConfig()
	if err != nil {
		slog.Error("failed to get configuration", "error", err)
		os.Exit(1)
	}

	database, err := db.Connect(config.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(database,
		&journal.MoodRecord{},
		&journal.DiaryEntry{},
		&account.User{},
		&session.RefreshToken{},
		&account.PasswordResetToken{},
		&account.Invite{},
	)
	if err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	e := server.NewEchoServer(config.Server)

	frontend.RegisterRoutes(e, config.Frontend)

	sessionRepository := session.NewRepository(database)
	sessionService := session.NewService(sessionRepository, config.Session)

	accountRepository := account.NewRepository(database)
	accountService := account.NewService(accountRepository, sessionService, config.Account)
	accountHandler := account.NewHandler(accountService, sessionService, config.Account, config.Session)
	adminService := admin.NewService(config.Admin.Password, sessionService)
	adminHandler := admin.NewHandler(accountService, adminService, config.Admin, config.Session)

	validateUser := func(ctx context.Context, userID uint) error {
		_, err := accountService.GetUserByID(ctx, userID)
		return err
	}

	userJWTMiddleware := session.UserJWTMiddleware(sessionService, validateUser)
	adminJWTMiddleware := session.AdminJWTMiddleware(sessionService)

	account.RegisterRoutes(e, accountHandler, userJWTMiddleware)
	admin.RegisterRoutes(e, adminHandler, adminJWTMiddleware)

	journalRepository := journal.NewRepository(database)
	journalService := journal.NewService(journalRepository)
	journalHandler := journal.NewHandler(journalService)

	journal.RegisterRoutes(e, journalHandler, userJWTMiddleware)

	return &App{
		sessionService: sessionService,
		echo:           e,
		config:         config,
	}
}

func (a *App) Start() {
	if err := a.sessionService.DeleteExpiredRefreshTokens(context.Background()); err != nil {
		slog.Error("failed to delete expired refresh tokens", "error", err)
		os.Exit(1)
	}

	if err := server.StartServer(a.echo, a.config.Server); err != nil {
		slog.Error("server stopped with an error", "error", err)
		os.Exit(1)
	}
}
