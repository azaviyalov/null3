package auth

import (
	"context"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitModule(e *echo.Echo, db *gorm.DB, config Config, stubUserConfig StubUserConfig) *Module {
	repo := NewRepository(db)
	ctx := context.Background()
	err := repo.DeleteExpiredRefreshTokens(ctx)
	if err != nil {
		logging.Error(ctx, "failed to delete expired refresh tokens, continuing", "error", err)
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
