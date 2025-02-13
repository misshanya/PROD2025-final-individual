package app

import (
	"context"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type CampaignService struct {
	repo repository.CampaignRepository
}

func NewCampaignService(repo repository.CampaignRepository) *CampaignService {
	return &CampaignService{
		repo: repo,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, advertiserID uuid.UUID, campaignRequest *domain.CampaignRequest) (*domain.Campaign, error) {
	campaign, err := s.repo.CreateCampaign(ctx, advertiserID, campaignRequest)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *CampaignService) GetCampaignsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID, size, page int) ([]domain.Campaign, error) {
	offset := size * page
	campaigns, err := s.repo.GetCampaignsByAdvertiserID(ctx, advertiserID, size, offset)
	if err != nil {
		return []domain.Campaign{}, err
	}
	return campaigns, nil
}