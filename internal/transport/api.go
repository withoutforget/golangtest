package transport

import (
	"context"
	"gotest/internal/database"
	"gotest/internal/repository"

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
	tx, err := database.NewTransaction(
		api.db, nil,
	)
	if err != nil {
		return c.JSON(fiber.Map{"error1": err.Error()})
	}
	defer tx.Close()

	r := repository.NewLogRepository(tx)

	_, err = r.AddLog(repository.NewLogModel{Raw: c.Query("r"), Level: c.Query("l")})

	if err != nil {
		tx.Rollback()
		return c.JSON(fiber.Map{"error2": err.Error()})
	}

	data, err := r.GetLogs()

	if err != nil {
		tx.Rollback()
		return c.JSON(fiber.Map{"error3": err.Error()})
	}

	tx.Commit()

	return c.JSON(
		fiber.Map{
			"data": data,
		})

}

func (s *Server) setupAPI() {
	api := NewAPI()
	s.fiber.Get("/", api.index_handler)
}
