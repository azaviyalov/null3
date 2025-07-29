package auth

import "github.com/labstack/echo/v4"

func RegisterAuthRoutes(e *echo.Echo) {
	e.POST("/api/auth/login", LoginHandler)
	e.GET("/api/auth/me", MeHandler, JWTMiddleware())
}
