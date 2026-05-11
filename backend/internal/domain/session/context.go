package session

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const echoActorKey = "internal/actor"

type Actor struct {
	UserID  uint
	IsAdmin bool
}

func GetActor(c echo.Context) (*Actor, error) {
	actor := c.Get(echoActorKey)
	if actor == nil {
		return nil, fmt.Errorf("actor is nil: %w", ErrActorNotAuthenticated)
	}

	value, ok := actor.(*Actor)
	if !ok {
		return nil, fmt.Errorf("invalid actor type: expected *session.Actor, got %T: %w", actor, ErrActorInvalidType)
	}

	return value, nil
}

func setActor(c echo.Context, actor *Actor) {
	c.Set(echoActorKey, actor)
}
