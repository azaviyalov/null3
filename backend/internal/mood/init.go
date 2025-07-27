package mood

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitModule(e *echo.Echo, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	RegisterRoutes(e, handler)
}
