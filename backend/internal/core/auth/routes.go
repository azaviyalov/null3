package auth

import "github.com/labstack/echo/v4"

func RegisterAuthRoutes(e *echo.Echo, config Config, stubUserConfig StubUserConfig, jwt echo.MiddlewareFunc) {
	e.POST("/api/auth/login", LoginHandler(config, stubUserConfig))
	e.POST("/api/auth/logout", LogoutHandler(config))
	e.GET("/api/auth/me", MeHandler, jwt)
}
