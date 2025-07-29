package auth

import (
	"github.com/labstack/echo/v4"
)

func InitModule(e *echo.Echo) {
	RegisterAuthRoutes(e)
}
