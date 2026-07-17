package app

import (
	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/frontend"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/admin"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
)

type Config struct {
	Admin    admin.Config
	Account  account.Config
	DB       db.Config
	Frontend frontend.Config
	Session  session.Config
	Server   server.Config
}

func GetConfig() (Config, error) {
	sessionConfig, err := session.GetConfig()
	if err != nil {
		return Config{}, err
	}

	accountConfig, err := account.GetConfig()
	if err != nil {
		return Config{}, err
	}

	dbConfig := db.GetConfig()

	frontendConfig, err := frontend.GetConfig()
	if err != nil {
		return Config{}, err
	}

	serverConfig, err := server.GetConfig()
	if err != nil {
		return Config{}, err
	}

	adminConfig, err := admin.GetConfig()
	if err != nil {
		return Config{}, err
	}
	adminConfig.FrontendURL = serverConfig.FrontendURL
	accountConfig.FrontendURL = serverConfig.FrontendURL

	return Config{
		Admin:    adminConfig,
		Account:  accountConfig,
		DB:       dbConfig,
		Frontend: frontendConfig,
		Session:  sessionConfig,
		Server:   serverConfig,
	}, nil
}
