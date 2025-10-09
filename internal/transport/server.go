package transport

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	fiber *fiber.App
	log   *slog.Logger
}

func NewServer(fiber *fiber.App,
	log *slog.Logger,
) *Server {
	return &Server{fiber: fiber, log: log}
}

func (s *Server) Run() {
	s.fiber.Use(logger.New())

	s.setupAPI()

	err := s.fiber.Listen(
		"localhost:80",
	)
	if err != nil {
		s.log.Error("Error during listen", slog.String("error", err.Error()))
		return
	}
	s.log.Info("Server stopped")
}

func (s *Server) Stop() error {
	return s.fiber.Shutdown()
}
