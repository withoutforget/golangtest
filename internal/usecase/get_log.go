package usecase

import "gotest/internal/repository"

type GetLogUsecase struct {
	log_repo *repository.LogRepository
}

func NewGetLogUsecase(lr *repository.LogRepository) *GetLogUsecase {
	return &GetLogUsecase{log_repo: lr}
}

func (u *GetLogUsecase) Run() ([]repository.LogModel, error) {
	logs, err := u.log_repo.GetLogs()
	return logs, err
}
