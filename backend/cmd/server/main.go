// Package main is the entry point for the IOrder API server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinhplus/IOrder/backend/internal/config"
	"github.com/dinhplus/IOrder/backend/internal/db"
	"github.com/dinhplus/IOrder/backend/internal/router"
)

func main() {
	cfg, cfgErr := config.Load()

	// Configure structured logger before anything else.
	setupLogger(cfg)

	if cfgErr != nil {
		slog.Warn("configuration warning", "error", cfgErr)
	}

	// Connect to database — server continues even if DB is temporarily unavailable.
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Warn("database unavailable at startup — running without DB", "error", err)
	}

	r := router.New(cfg, database)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background.
	go func() {
		slog.Info("server starting", "addr", srv.Addr, "env", cfg.Env)
		if serveErr := srv.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.Error("server error", "error", serveErr)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		slog.Error("graceful shutdown failed", "error", shutdownErr)
	}

	if database != nil {
		if closeErr := database.Close(); closeErr != nil {
			slog.Error("error closing database", "error", closeErr)
		}
	}

	slog.Info("server stopped")
}

func setupLogger(cfg *config.Config) {
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if cfg.Env == "production" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(h))
}
