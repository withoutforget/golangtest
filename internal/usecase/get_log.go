package usecase

import (
	"gotest/internal/repository"
	"time"
)

type GetLogUsecase struct {
	log_repo *repository.LogRepository
}

func NewGetLogUsecase(lr *repository.LogRepository) *GetLogUsecase {
	return &GetLogUsecase{log_repo: lr}
}

type GetLogUsecaseResponse struct {
	Logs []repository.LogModel `json:"logs"`
}

func (u *GetLogUsecase) Run(since *time.Time,
	before *time.Time,
	level []string,
	source *string,
	request_id *string,
	logger_name *string,
	group_asc bool) (GetLogUsecaseResponse, error) {
	logs, err := u.log_repo.GetLogs(since,
		before,
		level,
		source,
		request_id,
		logger_name,
		group_asc)
	return GetLogUsecaseResponse{Logs: logs}, err
}
