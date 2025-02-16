package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type CampaignService struct {
	repo          repository.CampaignRepository
	timeRepo      repository.TimeRepository
	openAIService domain.MLService
}

func NewCampaignService(repo repository.CampaignRepository, timeRepo repository.TimeRepository, openAIService domain.MLService) *CampaignService {
	return &CampaignService{
		repo:          repo,
		timeRepo:      timeRepo,
		openAIService: openAIService,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, advertiserID uuid.UUID, campaignRequest *domain.CampaignRequest) (*domain.Campaign, error) {
	// Combine ad title and ad text into one string to check both at once
	allText := fmt.Sprintf("Название: %s; Описание: %s", campaignRequest.AdTitle, campaignRequest.AdText)
	isAllowed, err := s.openAIService.ValidateAdText(ctx, allText)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		return nil, domain.ErrModerationNotPassed
	}
	currentDate, err := s.timeRepo.GetCurrentDate(ctx)
	if err != nil {
		return nil, err
	}
	campaign, err := s.repo.CreateCampaign(ctx, advertiserID, campaignRequest, *currentDate)
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

func (s *CampaignService) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*domain.Campaign, error) {
	campaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, campaignID uuid.UUID, campaignUpdate domain.CampaignUpdateRequest) (*domain.Campaign, error) {
	// Combine ad title and ad text into one string to check both at once
	allText := fmt.Sprintf("Название: %s; Описание: %s", campaignUpdate.AdTitle, campaignUpdate.AdText)
	isAllowed, err := s.openAIService.ValidateAdText(ctx, allText)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		return nil, domain.ErrModerationNotPassed
	}
	currentDate, err := s.timeRepo.GetCurrentDate(ctx)
	if err != nil {
		return nil, err
	}
	campaign, err := s.repo.UpdateCampaign(ctx, campaignID, campaignUpdate, *currentDate)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, campaignID uuid.UUID) error {
	err := s.repo.DeleteCampaign(ctx, campaignID)
	return err
}
