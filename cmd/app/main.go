package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/middleware/mwLogger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := setupLogger(cfg.Env)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("starting url-shortener", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	logger.Debug("storage initialized")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//TODO: start server

}

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, fmt.Errorf("uknown env mode: %s", env)
	}

	return logger, nil
}
