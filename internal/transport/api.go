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

func (api *API) append_log_handler(c *fiber.Ctx) error {
	raw := c.Query("raw")
	level := c.Query("level")

	isParamsValid, ret := validateParam(c, raw, level)
	if !isParamsValid {
		return ret
	}

	return withTx(c, api.db, func(tx database.TransactionManager) (any, error) {
		r := repository.NewLogRepository(tx)
		return usecase.NewAppendLogUsecase(r).Run(raw, level)
	})
}

func (api *API) get_log_handler(c *fiber.Ctx) error {
	return withTx(c, api.db, func(tx database.TransactionManager) (any, error) {
		r := repository.NewLogRepository(tx)
		return usecase.NewGetLogUsecase(r).Run()
	})
}

func (s *Server) setupAPI() {
	api := NewAPI()
	s.fiber.Get("/append", api.append_log_handler)
	s.fiber.Get("/get", api.get_log_handler)

}
