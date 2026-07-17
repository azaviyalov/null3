package session

import "github.com/labstack/echo/v4"

const echoUserIDKey = "internal/user-id"

func GetUserID(c echo.Context) uint {
	return c.Get(echoUserIDKey).(uint)
}

func setUserID(c echo.Context, userID uint) {
	c.Set(echoUserIDKey, userID)
}
