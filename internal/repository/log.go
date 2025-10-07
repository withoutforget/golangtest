package repository

import (
	"database/sql"
	"gotest/internal/database"
)

type LogModel struct {
	ID    uint64 `json:"id"`
	Raw   string `json:"raw"`
	Level string `json:"level"`
}

type NewLogModel struct {
	Raw   string
	Level string
}

type LogRepository struct {
	tx *sql.Tx
}

func NewLogRepository(tx database.TransactionManager) *LogRepository {
	return &LogRepository{tx: tx.Transaction()}
}

func (r *LogRepository) AddLog(model NewLogModel) (uint64, error) {
	rows, err := r.tx.Query(
		"INSERT INTO log (raw, level) VALUES ($1, $2) RETURNING id;",
		model.Raw,
		model.Level,
	)

	if err != nil {
		return 0, err
	}
	var tmp uint64
	for rows.Next() {
		err := rows.Scan(&tmp)
		if err != nil {
			return 0, err
		}
	}
	return tmp, nil

}

func (r *LogRepository) GetLogs() ([]LogModel, error) {
	rows, err := r.tx.Query("SELECT * FROM log;")
	if err != nil {
		return nil, err
	}
	res := make([]LogModel, 0)
	for rows.Next() {
		var id uint64
		var raw string
		var level string
		err := rows.Scan(&id, &raw, &level)
		if err != nil {
			return nil, err
		}
		res = append(res, LogModel{ID: id, Raw: raw, Level: level})
	}
	return res, nil
}
