package main

import (
	"context"
	"gotest/internal/logging"
	"gotest/internal/transport"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log, err := logging.InitLogging()
	if err != nil {
		panic(err.Error())
	}

	log.Info("Logger initialized")

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	fiber := fiber.New()

	srv := transport.NewServer(fiber, log)

	go srv.Run()

	GracefulShutdown(ctx, srv, log)
}

func GracefulShutdown(ctx context.Context, srv *transport.Server, log *slog.Logger) {
	<-ctx.Done()

	err := srv.Stop()
	if err != nil {
		log.Error("Error while graceful shutdown", slog.String("Error", err.Error()))
	}

	log.Info("Graceful shutdown...")
}
