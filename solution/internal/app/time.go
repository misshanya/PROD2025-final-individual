package app

import (
	"context"

	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type TimeService struct {
	repo repository.TimeRepository
}

func NewTimeService(repo repository.TimeRepository) *TimeService {
	return &TimeService{
		repo: repo,
	}
}

func (s *TimeService) SetCurrentDate(ctx context.Context, newDate int) error {
	return s.repo.SetCurrentDate(ctx, newDate)
}
