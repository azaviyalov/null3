package server

import (
	"os"

	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/go-playground/validator/v10"
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

	e.Validator = newCustomValidator()

	auth.RegisterAuthRoutes(e)

	// Add a health check endpoint for frontend to check backend status
	e.GET("/api/healthz", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"status": "ok"})
	})

	return e
}

type customValidator struct {
	v *validator.Validate
}

func newCustomValidator() *customValidator {
	return &customValidator{v: validator.New(validator.WithRequiredStructEnabled())}
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}
