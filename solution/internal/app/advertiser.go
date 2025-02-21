package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type AdvertiserService struct {
	repo     repository.AdvertiserRepository
	userRepo repository.UserRepository
}

func NewAdvertiserService(repo repository.AdvertiserRepository,
	userRepo repository.UserRepository) *AdvertiserService {
	return &AdvertiserService{
		repo:     repo,
		userRepo: userRepo,
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
	_, err := s.repo.GetByID(ctx, score.AdvertiserID)
	if err != nil {
		return nil, err
	}
	_, err = s.userRepo.GetByID(ctx, score.ClientID)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetMLScore(ctx, score.ClientID, score.AdvertiserID); err == pgx.ErrNoRows {
		err := s.repo.CreateMLScore(ctx, score)
		if err != nil {
			return nil, err
		}
		return score, nil
	}
	s.repo.UpdateMLScore(ctx, score)
	return score, nil
}
