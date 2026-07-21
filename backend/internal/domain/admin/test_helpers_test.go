package admin_test

import (
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/admin"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type adminTestEnvironment struct {
	database       *gorm.DB
	echo           *echo.Echo
	sessionService *session.Service
}

func newAdminTestEnvironment(t *testing.T) *adminTestEnvironment {
	t.Helper()

	database := testutil.NewDatabase(t, "admin.sqlite")

	sessionConfig := session.Config{
		JWTSecret:              testJWTSecret,
		JWTExpiration:          time.Hour,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
		SecureCookies:          true,
	}
	sessionService := session.NewService(session.NewRepository(database), sessionConfig)
	accountService := account.NewService(account.NewRepository(database), sessionService, account.Config{
		PasswordResetTokenExpiration: time.Hour,
		FrontendURL:                  "https://journal.example",
	})
	adminService := admin.NewService(testAdminPassword, sessionService)

	testutil.DiscardLogs(t)

	e := server.NewEchoServer(server.Config{})
	handler := admin.NewHandler(accountService, adminService, admin.Config{
		FrontendURL: "https://journal.example",
		Password:    testAdminPassword,
	}, sessionConfig)
	admin.RegisterRoutes(e, handler, session.AdminJWTMiddleware(sessionService))

	return &adminTestEnvironment{
		database:       database,
		echo:           e,
		sessionService: sessionService,
	}
}
