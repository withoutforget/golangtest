package database

import (
	"database/sql"
)

type Migration struct {
	tx      *sql.Tx
	queries []string
}

func NewMigration(tx TransactionManager, q ...string) *Migration {
	qr := make([]string, 0)
	qr = append(qr, q...)
	return &Migration{tx: tx.Transaction(), queries: qr}
}

func (m *Migration) Run() error {
	for _, q := range m.queries {
		_, err := m.tx.Query(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func Migrate(tx TransactionManager) error {
	m1 := NewMigration(
		tx, "CREATE TABLE log (id BIGINT PRIMARY KEY, raw TEXT, level VARCHAR(50));",
	)

	err := m1.Run()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
