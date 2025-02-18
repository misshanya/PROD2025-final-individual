package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type AdsService struct {
	repo         repository.AdsRepository
	userRepo     repository.UserRepository
	campaignRepo repository.CampaignRepository
}

func NewAdsService(repo repository.AdsRepository, userRepo repository.UserRepository, campaignRepo repository.CampaignRepository) *AdsService {
	return &AdsService{
		repo:         repo,
		userRepo:     userRepo,
		campaignRepo: campaignRepo,
	}
}

func (s *AdsService) Click(ctx context.Context, adId, clientId uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, clientId)
	if err == pgx.ErrNoRows {
		return domain.ErrUserNotFound
	} else if err != nil {
		return err
	}

	// Check if ad exists
	_, err = s.campaignRepo.GetCampaignByID(ctx, adId)
	if err == pgx.ErrNoRows {
		return domain.ErrAdNotFound
	} else if err != nil {
		return err
	}

	return s.repo.Click(ctx, adId, clientId)
}
