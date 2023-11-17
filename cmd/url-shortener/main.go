package main

import (
	"api-shorter/internal/config"
	mwLogger "api-shorter/internal/http-server/middleware/logger"
	"api-shorter/internal/lib/logger/sl"
	"api-shorter/internal/storage/pgsql"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting url shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := pgsql.New()

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	err = storage.DeleteURL("vk")
	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
