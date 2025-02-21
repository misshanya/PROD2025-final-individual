package app

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *AdvertiserService) CreateUpdateAdvertisers(ctx context.Context, advertisers []*domain.Advertiser) ([]*domain.Advertiser, error) {
	newAdvertisers, err := s.repo.CreateUpdateAdvertisers(ctx, advertisers)
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

func (s *AdvertiserService) CreateUpdateMLScore(ctx context.Context, score *domain.MLScore) (*domain.MLScore, error) {
	if _, err := s.repo.GetMLScore(ctx, score.ClientID, score.AdvertiserID); err == pgx.ErrNoRows {
		log.Println("Creating score")
		err := s.repo.CreateMLScore(ctx, score)
		if err != nil {
			return nil, err
		}
		return score, nil
	}
	log.Println("Updating score")
	s.repo.UpdateMLScore(ctx, score)
	return score, nil
}
