package repository

import (
	"database/sql"
	"gotest/internal/database"
	"time"

	"github.com/lib/pq"
)

type LogModel struct {
	ID         uint64    `json:"id"`
	Raw        string    `json:"raw"`
	Level      string    `json:"level"`
	CreatedAt  time.Time `json:"created_at"`
	Source     string    `json:"source"`
	RequestID  string    `json:"request_id"`
	LoggerName string    `json:"logger_name"`
}

type NewLogModel struct {
	Raw        string    `json:"raw"`
	Level      string    `json:"level"`
	CreatedAt  time.Time `json:"created_at"`
	Source     string    `json:"source"`
	RequestID  string    `json:"request_id"`
	LoggerName string    `json:"logger_name"`
}

type LogRepository struct {
	tx *sql.Tx
}

func NewLogRepository(tx database.TransactionManager) *LogRepository {
	return &LogRepository{tx: tx.Transaction()}
}

func (r *LogRepository) AddLog(model NewLogModel) (uint64, error) {

	rows, err := r.tx.Query(
		`INSERT INTO log
		 (raw, level,
		 created_at, 
		 source, 
		 request_id, 
		 logger_name)
		  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		model.Raw,
		model.Level,
		model.CreatedAt,
		model.Source,
		model.RequestID,
		model.LoggerName,
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

func (r *LogRepository) GetLogs(
	since *time.Time,
	before *time.Time,
	level []string,
	source *string,
	request_id *string,
	logger_name *string) ([]LogModel, error) {
	rows, err := r.tx.Query(`SELECT 
	id, 
	raw, 
	level, 
	created_at, 
	source,
	request_id, 
	logger_name 
	FROM log WHERE 
		($1::timestamptz IS NULL OR created_at >= $1)
		AND ($2::timestamptz IS NULL OR created_at <= $2)
		AND ($3::TEXT[] IS NULL OR level = ANY($3))
		AND ($4::TEXT IS NULL OR source = $4)
		AND ($5::TEXT IS NULL OR request_id = $5)
		AND ($6::TEXT IS NULL OR logger_name = $6)
	;
	`,
		since,
		before,
		pq.Array(level),
		source,
		request_id,
		logger_name)
	if err != nil {
		return nil, err
	}

	res := make([]LogModel, 0)
	for rows.Next() {
		var id uint64
		var raw string
		var level string
		var created_at time.Time
		var source string
		var request_id string
		var logger_name string
		err := rows.Scan(&id, &raw, &level, &created_at, &source, &request_id, &logger_name)
		if err != nil {
			return nil, err
		}
		res = append(res, LogModel{ID: id,
			Raw:        raw,
			Level:      level,
			CreatedAt:  created_at,
			Source:     source,
			RequestID:  request_id,
			LoggerName: logger_name})
	}
	return res, nil
}
