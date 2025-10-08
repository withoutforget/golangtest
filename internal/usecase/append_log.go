package usecase

import "gotest/internal/repository"

type AppendLogUsecase struct {
	log_repo *repository.LogRepository
}

func NewLogUsecase(lr *repository.LogRepository) *AppendLogUsecase {
	return &AppendLogUsecase{log_repo: lr}
}

func (u *AppendLogUsecase) Run(raw string, level string) (uint64, error) {
	id, err := u.log_repo.AddLog(repository.NewLogModel{Raw: raw, Level: level})
	return id, err
}
