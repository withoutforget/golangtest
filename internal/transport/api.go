package transport

import (
	"context"
	"gotest/internal/database"
	"gotest/internal/repository"
	"gotest/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type API struct {
	db  *database.Database
	ctx context.Context
}

func NewAPI() *API {
	ctx := context.Background()
	db, err := database.NewDatabase(ctx)

	if err != nil {
		panic(err.Error())
	}

	return &API{db: db, ctx: ctx}
}

func (api *API) index_handler(c *fiber.Ctx) error {
	raw := c.Query("raw")
	level := c.Query("level")

	if raw == "" || level == "" {
		return c.Status(400).JSON(fiber.Map{"error": "incorrect input"})
	}
	tx, err := database.NewTransaction(api.db, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer tx.Close()
	r := repository.NewLogRepository(tx)

	u := usecase.NewLogUsecase(r)

	id, err := u.Run(raw, level)

	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tx.Commit()

	return c.Status(200).JSON(fiber.Map{"id": id})
}

func (s *Server) setupAPI() {
	api := NewAPI()
	s.fiber.Get("/", api.index_handler)
}
