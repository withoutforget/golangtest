package usecase

import "gotest/internal/repository"

type AppendLogUsecase struct {
	log_repo *repository.LogRepository
}

func NewAppendLogUsecase(lr *repository.LogRepository) *AppendLogUsecase {
	return &AppendLogUsecase{log_repo: lr}
}

type AppendLogUsecaseResponse struct {
	ID uint64 `json:"id"`
}

func (u *AppendLogUsecase) Run(raw string, level string) (AppendLogUsecaseResponse, error) {
	id, err := u.log_repo.AddLog(repository.NewLogModel{Raw: raw, Level: level})
	return AppendLogUsecaseResponse{ID: id}, err
}
