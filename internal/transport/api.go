package transport

import (
	"context"
	"gotest/internal/database"

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

	tx, err := database.NewTransaction(db, nil)
	if err != nil {
		panic(err.Error())
	}
	defer tx.Commit()
	defer tx.Close()

	err = database.Migrate(tx)
	if err != nil {
		panic(err.Error())
	}

	return &API{db: db, ctx: ctx}
}

func (api *API) index_handler(c *fiber.Ctx) error {
	tr, err := database.NewTransaction(api.db, nil)
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	defer tr.Close()
	defer tr.Rollback()

	q1 := "CREATE TABLE users (id INT PRIMARY KEY);INSERT INTO users (id) VALUES (1), (2), (3), (4), (5), (6), (7), (8), (9), (10);"
	tr.Transaction().Query(q1)

	rows, err := tr.Transaction().Query("SELECT * FROM users;")
	if err != nil {
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	array := make([]int, 0)
	for rows.Next() {
		var tmp int

		if err := rows.Scan(&tmp); err != nil {
			return c.JSON(fiber.Map{"error": err.Error()})
		}
		array = append(array, tmp)
	}
	return c.JSON(fiber.Map{"result": array})
}

func (s *Server) setupAPI() {
	api := NewAPI()
	s.fiber.Get("/", api.index_handler)
}
