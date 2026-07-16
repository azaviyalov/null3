package app

import (
	"context"
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
		&journal.MoodEntry{},
		&journal.DiaryEntry{},
		&account.User{},
		&session.RefreshToken{},
		&account.PasswordResetToken{},
		&account.Invite{},
	)
	if err != nil {
		logging.Error(context.Background(), "database migration failed", "error", err)
		os.Exit(1)
	}

	e := server.NewEchoServer(config.Server)

	frontend.RegisterRoutes(e, config.Frontend)

	sessionRepository := session.NewRepository(database)
	sessionService := session.NewService(sessionRepository, config.Session)

	accountRepository := account.NewRepository(database)
	accountService := account.NewService(accountRepository, sessionService, config.Account)
	if err := accountService.SeedAdminUser(context.Background()); err != nil {
		logging.Error(context.Background(), "failed to seed admin user", "error", err)
		os.Exit(1)
	}
	accountHandler := account.NewHandler(accountService, sessionService, config.Account, config.Session)
	adminHandler := admin.NewHandler(accountService, sessionService, config.Admin, config.Session)

	resolveActor := func(ctx context.Context, userID uint) (*session.Actor, error) {
		user, err := accountService.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		return &session.Actor{
			UserID:  user.ID,
			IsAdmin: accountService.IsAdmin(user),
		}, nil
	}

	userJWTMiddleware := session.UserJWTMiddleware(sessionService, resolveActor)
	adminJWTMiddleware := session.AdminJWTMiddleware(sessionService, resolveActor)

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
		logging.Error(context.Background(), "failed to delete expired refresh tokens", "error", err)
		os.Exit(1)
	}

	if err := server.StartServer(a.echo, a.config.Server); err != nil {
		logging.Error(context.Background(), "server stopped with an error", "error", err)
		os.Exit(1)
	}
}
