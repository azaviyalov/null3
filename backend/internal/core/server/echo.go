package server

import (
	"os"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	if os.Getenv("ENABLE_CORS") == "true" {
		e.Use(middleware.CORS())
	}
	e.Use(logging.RequestLogger())
	e.Use(middleware.Recover())

	return e
}
