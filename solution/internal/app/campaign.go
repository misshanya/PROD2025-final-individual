package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type CampaignService struct {
	repo            repository.CampaignRepository
	advertiserRepo  repository.AdvertiserRepository
	timeRepo        repository.TimeRepository
	openAIService   domain.MLService
	mlRepository    repository.MLRepository
	fileRepo        repository.FileRepository
	minioPublicHost string
}

func NewCampaignService(repo repository.CampaignRepository,
	advertiserRepo repository.AdvertiserRepository,
	timeRepo repository.TimeRepository,
	openAIService domain.MLService,
	mlRepository repository.MLRepository,
	fileRepo repository.FileRepository,
	minioPublicHost string) *CampaignService {
	return &CampaignService{
		repo:            repo,
		advertiserRepo:  advertiserRepo,
		timeRepo:        timeRepo,
		openAIService:   openAIService,
		mlRepository:    mlRepository,
		fileRepo:        fileRepo,
		minioPublicHost: minioPublicHost,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, advertiserID uuid.UUID, campaignRequest *domain.CampaignRequest) (*domain.Campaign, error) {
	isModerated, err := s.checkModeration(ctx)
	if err != nil {
		return nil, err
	}
	if isModerated {
		// Combine ad title and ad text into one string to check both at once
		allText := fmt.Sprintf("Название: %s; Описание: %s", campaignRequest.AdTitle, campaignRequest.AdText)
		isAllowed, err := s.openAIService.ValidateAdText(ctx, allText)
		if err != nil {
			return nil, err
		}
		if !isAllowed {
			return nil, domain.ErrModerationNotPassed
		}
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

func (s *CampaignService) SetCampaignPicture(ctx context.Context,
	advertiserID, campaignID uuid.UUID,
	fileName string,
	fileContent []byte) error {
	_, err := s.advertiserRepo.GetByID(ctx, advertiserID)
	if err == pgx.ErrNoRows {
		return domain.ErrAdvertiserNotFound
	}
	_, err = s.repo.GetCampaignByID(ctx, campaignID)
	if err == pgx.ErrNoRows {
		return domain.ErrAdNotFound
	}

	// Generate file key
	id := uuid.New().String()
	ext := filepath.Ext(fileName)
	fileKey := id + ext

	err = s.fileRepo.UploadFile(ctx, fileKey, fileContent)
	if err != nil {
		return err
	}

	return s.repo.SetCampaignPicture(ctx, campaignID, fileKey)
}

func (s *CampaignService) GetCampaignsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID, size, page int) ([]domain.Campaign, error) {
	offset := size * page
	campaigns, err := s.repo.GetCampaignsByAdvertiserID(ctx, advertiserID, size, offset)
	if err != nil {
		return []domain.Campaign{}, err
	}
	for i := range campaigns {
		picID, err := s.repo.GetCampaignPicID(ctx, campaigns[i].ID)
		if err != nil || picID == "" {
			continue
		}
		picURL, err := s.fileRepo.GetFileLink(ctx, picID, s.minioPublicHost)
		if err != nil || picURL == "" {
			continue
		}
		campaigns[i].PicURL = &picURL
	}
	return campaigns, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*domain.Campaign, error) {
	campaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	picID, err := s.repo.GetCampaignPicID(ctx, campaignID)
	if err != nil || picID == "" {
		return campaign, nil
	}
	picURL, err := s.fileRepo.GetFileLink(ctx, picID, s.minioPublicHost)
	if err != nil || picURL == "" {
		return campaign, nil
	}
	campaign.PicURL = &picURL
	return campaign, nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, campaignID uuid.UUID, campaignUpdate domain.CampaignUpdateRequest) (*domain.Campaign, error) {
	isModerated, err := s.checkModeration(ctx)
	if err != nil {
		return nil, err
	}
	if isModerated {
		// Combine ad title and ad text into one string to check both at once
		allText := fmt.Sprintf("Название: %s; Описание: %s", campaignUpdate.AdTitle, campaignUpdate.AdText)
		isAllowed, err := s.openAIService.ValidateAdText(ctx, allText)
		if err != nil {
			return nil, err
		}
		if !isAllowed {
			return nil, domain.ErrModerationNotPassed
		}
	}
	currentDate, err := s.timeRepo.GetCurrentDate(ctx)
	if err != nil {
		return nil, err
	}
	campaign, err := s.repo.UpdateCampaign(ctx, campaignID, campaignUpdate, *currentDate)
	if err != nil {
		return nil, err
	}
	picID, err := s.repo.GetCampaignPicID(ctx, campaignID)
	if err != nil || picID == "" {
		return campaign, nil
	}
	picURL, err := s.fileRepo.GetFileLink(ctx, picID, s.minioPublicHost)
	if err != nil || picURL == "" {
		return campaign, nil
	}
	campaign.PicURL = &picURL
	return campaign, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, campaignID uuid.UUID) error {
	err := s.repo.DeleteCampaign(ctx, campaignID)
	return err
}

func (s *CampaignService) GenerateAdText(ctx context.Context, advertiserName string, adTitle string) (string, error) {
	adText, err := s.openAIService.GenerateAdText(ctx, advertiserName, adTitle)
	return adText, err
}

func (s *CampaignService) SwitchModeration(ctx context.Context) (bool, error) {
	return s.mlRepository.SwitchModeration(ctx)
}

func (s *CampaignService) checkModeration(ctx context.Context) (bool, error) {
	return s.mlRepository.CheckModeration(ctx)
}
