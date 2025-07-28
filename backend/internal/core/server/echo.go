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

	// Enable cors only if the frontend is not served from the same server.
	if os.Getenv("ENABLE_FRONTEND_DIST") != "true" {
		e.Use(middleware.CORS())
	}
	e.Use(logging.RequestLogger())
	e.Use(middleware.Recover())

	return e
}
