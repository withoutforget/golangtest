package transport

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gotest/internal/database"
)

func withTx(
	c *fiber.Ctx,
	db *database.Database,
	fn func(tx database.TransactionManager) (any, error),
) error {
	tx, err := database.NewTransaction(db, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer tx.Close()

	res, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tx.Commit()

	return c.Status(200).JSON(res)
}

func validateParam(c *fiber.Ctx, in ...string) (bool, error) {
	for _, v := range in {
		slog.Default().Info("1")
		if len(v) == 0 {
			slog.Default().Error("2")
			return false, c.Status(400).JSON(fiber.Map{"error": "incorrect input"})
		}
	}
	return true, nil
}
