package transport

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

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

type AppendLogRequest struct {
	Raw        string    `json:"raw"`
	Level      string    `json:"level"`
	CreatedAt  time.Time `json:"created_at"`
	Source     *string   `json:"source"`
	RequestID  *string   `json:"request_id"`
	LoggerName *string   `json:"logger_name"`
}

func (api *API) append_log_handler(c *fiber.Ctx) error {
	var request AppendLogRequest
	request_raw := c.Body()
	err := json.Unmarshal(request_raw, &request)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return withTx(c, api.db, func(tx database.TransactionManager) (any, error) {
		r := repository.NewLogRepository(tx)
		return usecase.NewAppendLogUsecase(r).Run(repository.NewLogModel{
			Raw:        request.Raw,
			Level:      request.Level,
			CreatedAt:  request.CreatedAt,
			Source:     request.Source,
			RequestID:  request.RequestID,
			LoggerName: request.LoggerName,
		})
	})
}

type GetLogHandlerRequest struct {
	Since       *time.Time `json:"since"`
	Before      *time.Time `json:"before"`
	Level       []string   `json:"level"`
	Source      *string    `json:"source"`
	RequesID    *string    `json:"request_id"`
	Logger_name *string    `json:"logger_name"`
	Order       bool       `json:"order"`
}

func (api *API) get_log_handler(c *fiber.Ctx) error {

	var since *time.Time
	var before *time.Time
	var level []string
	var source *string
	var request_id *string
	var logger_name *string
	var order bool

	sinceq := c.Query("since")
	beforeq := c.Query("before")
	levelq := c.Query("level")
	sourceq := c.Query("source")
	request_idq := c.Query("request_id")
	logger_nameq := c.Query("logger_name")
	orderq := c.Query("order")

	if orderq == "" || orderq == "ASC" {
		order = true
	} else {
		order = false
	}

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

	data, err := json.Marshal(GetLogHandlerRequest{since,
		before,
		level,
		source,
		request_id,
		logger_name,
		order})

	if err != nil {
		panic("abc")
	}

	hasher := sha256.New()
	hasher.Write(data)
	res := hasher.Sum(nil)
	_ = base64.URLEncoding.EncodeToString(res)

	return withTx(c, api.db, func(tx database.TransactionManager) (any, error) {
		r := repository.NewLogRepository(tx)
		return usecase.NewGetLogUsecase(r).Run(
			since,
			before,
			level,
			source,
			request_id,
			logger_name,
			order,
		)
	})
}

func (s *Server) setupAPI() {
	api := NewAPI()
	s.fiber.Post("/append", api.append_log_handler)
	s.fiber.Get("/get", api.get_log_handler)
}
