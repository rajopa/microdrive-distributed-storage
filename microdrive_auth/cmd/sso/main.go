package main

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	grpcapp "microdrive_auth/internal/app/grpc"
	"microdrive_auth/internal/config"
	"microdrive_auth/internal/services/auth"
	"microdrive_auth/internal/storage/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	m, err := migrate.New(
		"file://"+"./migrations",
		cfg.StoragePath,
	)
	if err == nil {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Error("migration error", slog.Any("err", err))
			panic(err)
		}
	}

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, cfg.TokenTTL)

	application := grpcapp.New(log, authService, cfg.GRPC.Port)

	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	log.Info("application started", slog.Int("port", cfg.GRPC.Port))

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.Stop()
	log.Info("application stopped gracefully")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
