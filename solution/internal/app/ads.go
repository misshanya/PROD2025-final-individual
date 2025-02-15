package app

import (
	"context"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type AdsService struct {
	repo repository.AdsRepository
}

func NewAdsService(repo repository.AdsRepository) *AdsService {
	return &AdsService{
		repo: repo,
	}
}

func (s *AdsService) Click(ctx context.Context, adId, clientId uuid.UUID) error {
	return s.repo.Click(ctx, adId, clientId)
}