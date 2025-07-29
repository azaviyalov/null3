package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/auth"
	"github.com/azaviyalov/null3/backend/internal/core/db"
	"github.com/azaviyalov/null3/backend/internal/core/env"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/frontend"
	"github.com/azaviyalov/null3/backend/internal/mood"
)

func main() {
	env.Setup()
	logging.Setup()

	database := db.Connect()
	db.AutoMigrate(database)

	e := server.NewEchoServer()

	auth.InitModule(e)
	mood.InitModule(e, database)

	if os.Getenv("ENABLE_FRONTEND_DIST") == "true" {
		slog.Info("serving frontend dist", "API_URL", os.Getenv("API_URL"))
		frontend.InitModule(e)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		slog.Info("starting HTTP server", "port", os.Getenv("PORT"))
		if err := e.Start(":" + os.Getenv("PORT")); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("server start failed", "error", err)
				os.Exit(1)
			}
		}
	}()

	<-ctx.Done()
	slog.Info("received shutdown signal, shutting down server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped successfully")
}
