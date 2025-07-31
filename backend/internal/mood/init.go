package mood

import (
	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitModule(e *echo.Echo, db *gorm.DB, authModule *auth.Module) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	RegisterRoutes(e, handler, authModule.JWTMiddleware)
}
