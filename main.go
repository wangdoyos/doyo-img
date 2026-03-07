package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/wangdoyos/doyo-img/internal/cleanup"
	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/router"
	"github.com/wangdoyos/doyo-img/internal/service"
	"github.com/wangdoyos/doyo-img/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger
	level := slog.LevelInfo
	switch cfg.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	// Initialize storage
	var store storage.Storage
	switch cfg.Storage.Type {
	case "local":
		store, err = storage.NewLocalStorage(cfg.Storage.Local.DataDir)
		if err != nil {
			slog.Error("failed to initialize local storage", "error", err)
			os.Exit(1)
		}
		slog.Info("using local storage", "data_dir", cfg.Storage.Local.DataDir)
	case "s3":
		store, err = storage.NewS3Storage(&cfg.Storage.S3)
		if err != nil {
			slog.Error("failed to initialize S3 storage", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("unsupported storage type", "type", cfg.Storage.Type)
		os.Exit(1)
	}

	// Initialize service
	svc := service.NewImageService(store, cfg)

	// Setup Gin
	if cfg.Log.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	// Setup routes with embedded frontend
	frontendFS := FrontendFS()
	router.Setup(r, svc, cfg, frontendFS)

	// Start cleanup scheduler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cleanup.Start(ctx, store, &cfg.Cleanup)

	// Start server
	addr := cfg.Address()
	slog.Info("doyo-img server starting", "address", addr)

	go func() {
		if err := r.Run(addr); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")
	cancel()
}
