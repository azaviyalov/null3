package auth

import (
	"github.com/labstack/echo/v4"
)

func InitModule(e *echo.Echo, config Config, stubUserConfig StubUserConfig) *Module {
	jwt := JWTMiddleware(config)
	RegisterAuthRoutes(e, config, stubUserConfig, jwt)
	return &Module{
		JWTMiddleware: jwt,
	}
}

type Module struct {
	JWTMiddleware echo.MiddlewareFunc
}
