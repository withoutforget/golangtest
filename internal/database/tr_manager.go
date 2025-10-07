package database

import (
	"database/sql"
	"log/slog"
)

type TransactionManager interface {
	IsValid() bool
	Commit() error
	Rollback() error
	Close()
	Connection() *sql.Conn
	Transaction() *sql.Tx
}

type TransactionManagerImpl struct {
	connection  *sql.Conn
	transaction *sql.Tx
	isValid     bool
	hascommited bool
}

func NewTransaction(db *Database, opts *sql.TxOptions) (*TransactionManagerImpl, error) {
	log := slog.Default()

	ret := &TransactionManagerImpl{hascommited: false}
	conn, err := db.GetConnection()

	if err != nil {
		log.Error("Cannot create connection", slog.String("error", err.Error()))
		return nil, err
	}

	tx, err := conn.BeginTx(db.ctx, opts)
	if err != nil {
		conn_err := conn.Close()
		if conn_err != nil {
			log.Error("Cannot close connection", slog.String("error", err.Error()))
		}
		return nil, err
	}

	ret.connection = conn
	ret.transaction = tx
	return ret, nil
}

func (t *TransactionManagerImpl) IsValid() bool {
	return t.isValid
}

func (t *TransactionManagerImpl) Commit() error {
	return t.transaction.Commit()
}

func (t *TransactionManagerImpl) Rollback() error {
	return t.transaction.Rollback()
}

func (t *TransactionManagerImpl) Close() {
	err := t.connection.Close()
	if err != nil {
		slog.Default().Error("Cannot close connection", slog.String("error", err.Error()))
	}
}

func (t *TransactionManagerImpl) Connection() *sql.Conn {
	return t.connection
}

func (t *TransactionManagerImpl) Transaction() *sql.Tx {
	return t.transaction
}
