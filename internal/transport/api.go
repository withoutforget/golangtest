package transport

import (
	"context"
	"gotest/internal/database"
	"gotest/internal/repository"
	"gotest/internal/usecase"
	"strconv"
	"time"

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
	created_at := c.Query("created_at")
	source := c.Query("source")
	request_id := c.Query("request_id")
	logger_name := c.Query("logger_name")

	isParamsValid, ret := validateParam(c, raw, level, created_at)
	if !isParamsValid {
		return ret
	}

	created_at_timestamp, err := strconv.ParseInt(created_at, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return withTx(c, api.db, func(tx database.TransactionManager) (any, error) {
		r := repository.NewLogRepository(tx)
		return usecase.NewAppendLogUsecase(r).Run(repository.NewLogModel{
			Raw:        raw,
			Level:      level,
			CreatedAt:  time.Unix(created_at_timestamp, 0),
			Source:     source,
			RequestID:  request_id,
			LoggerName: logger_name,
		})
	})
}

func (api *API) get_log_handler(c *fiber.Ctx) error {
	var since *time.Time
	var before *time.Time
	var level []string
	var source *string
	var request_id *string
	var logger_name *string

	sinceq := c.Query("since")
	beforeq := c.Query("before")
	levelq := c.Query("level")
	sourceq := c.Query("source")
	request_idq := c.Query("request_id")
	logger_nameq := c.Query("logger_name")

	if sinceq != "" {
		v, err := strconv.ParseInt(sinceq, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		t := time.Unix(v, 0)
		since = &t
	}

	if beforeq != "" {
		v, err := strconv.ParseInt(beforeq, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		t := time.Unix(v, 0)
		before = &t
	}

	if levelq != "" {
		level = make([]string, 0)
		level = append(level, levelq)
	}
	if sourceq != "" {
		source = &sourceq
	}
	if request_idq != "" {
		request_id = &request_idq
	}
	if logger_nameq != "" {
		logger_name = &logger_nameq
	}

	return withTx(c, api.db, func(tx database.TransactionManager) (any, error) {
		r := repository.NewLogRepository(tx)
		return usecase.NewGetLogUsecase(r).Run(
			since,
			before,
			level,
			source,
			request_id,
			logger_name,
		)
	})
}

func (s *Server) setupAPI() {
	api := NewAPI()
	s.fiber.Get("/append", api.append_log_handler)
	s.fiber.Get("/get", api.get_log_handler)

}
