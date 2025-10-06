package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	pool *sql.DB
	ctx  context.Context
}

func NewDatabase(ctx context.Context) (*Database, error) {
	pool, err := getDatabase()
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool, ctx: ctx}, nil
}

func (d *Database) GetConnection() (*sql.Conn, error) {
	con, err := d.pool.Conn(d.ctx)
	return con, err
}

func getDatabase() (*sql.DB, error) {
	driver := "postgres"
	username := "postgres"
	password := "postgres"
	host := "localhost"
	port := "5432"
	database_name := "postgres"
	dsn := fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable", driver, username, password, host, port, database_name)
	con, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	con.SetMaxOpenConns(100)

	return con, nil
}
