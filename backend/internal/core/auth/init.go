package auth

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitModule(e *echo.Echo, db *gorm.DB, config Config, stubUserConfig StubUserConfig) *Module {
	repo := NewRepository(db)
	err := repo.DeleteExpiredRefreshTokens()
	if err != nil {
		slog.Error("failed to delete expired refresh tokens", "error", err)
	}
	service := NewService(repo, config, stubUserConfig)
	jwt := JWTMiddleware(config, service)
	handler := NewHandler(service, config, stubUserConfig)
	RegisterRoutes(e, handler, jwt)
	return &Module{
		JWTMiddleware: jwt,
	}
}

type Module struct {
	JWTMiddleware echo.MiddlewareFunc
}
