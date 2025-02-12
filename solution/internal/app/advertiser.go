package app

import (
	"context"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type AdvertiserService struct {
	repo repository.AdvertiserRepository
}

func NewAdvertiserService(repo repository.AdvertiserRepository) *AdvertiserService {
	return &AdvertiserService{
		repo: repo,
	}
}

func (s *AdvertiserService) CreateAdvertisers(ctx context.Context, advertisers []*domain.Advertiser) ([]*domain.Advertiser, error) {
	newAdvertisers, err := s.repo.Create(ctx, advertisers)
	if err != nil {
		return []*domain.Advertiser{}, err
	}

	return newAdvertisers, nil  
}

func (s *AdvertiserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Advertiser, error) {
	advertiser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return advertiser, nil
}