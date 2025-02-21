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
	// Check if advertiser exists
	_, err := s.advertiserRepo.GetByID(ctx, advertiserID)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrAdvertiserNotFound
	} else if err != nil {
		return nil, err
	}

	// Validate targeting
	if !validateTargeting(campaignRequest.Targeting) {
		return nil, domain.ErrBadRequest
	}

	if err := s.validateModeration(ctx, campaignRequest.AdTitle, campaignRequest.AdText); err != nil {
		return nil, err
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
	_, err := s.advertiserRepo.GetByID(ctx, advertiserID)
	if err != nil {
		return []domain.Campaign{}, domain.ErrAdvertiserNotFound
	}
	offset := size * page
	campaigns, err := s.repo.GetCampaignsByAdvertiserID(ctx, advertiserID, size, offset)
	if err != nil {
		return []domain.Campaign{}, err
	}
	for i := range campaigns {
		picURL, err := s.getPicURL(ctx, campaigns[i].ID)
		if err == nil {
			// Set campaign pic url
			campaigns[i].PicURL = &picURL
		}
	}
	return campaigns, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, advertiserID, campaignID uuid.UUID) (*domain.Campaign, error) {
	_, err := s.advertiserRepo.GetByID(ctx, advertiserID)
	if err != nil {
		return nil, domain.ErrAdvertiserNotFound
	}
	campaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrAdNotFound
	}
	picURL, err := s.getPicURL(ctx, campaignID)
	if err == nil {
		// Set campaign pic url
		campaign.PicURL = &picURL
	}
	return campaign, nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, advertiserID, campaignID uuid.UUID, campaignUpdate domain.CampaignUpdateRequest) (*domain.Campaign, error) {
	// Check if advertiser exists
	_, err := s.advertiserRepo.GetByID(ctx, advertiserID)
	if err != nil {
		return nil, domain.ErrAdvertiserNotFound
	}
	// Check if campaign exists
	_, err = s.repo.GetCampaignByID(ctx, campaignID)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrAdNotFound
	} else if err != nil {
		return nil, err
	}

	// Validate targeting
	if !validateTargeting(campaignUpdate.Targeting) {
		return nil, domain.ErrBadRequest
	}

	if err := s.validateModeration(ctx, campaignUpdate.AdTitle, campaignUpdate.AdText); err != nil {
		return nil, err
	}

	// Get current day
	currentDate, err := s.timeRepo.GetCurrentDate(ctx)
	if err != nil {
		return nil, err
	}

	campaign, err := s.repo.UpdateCampaign(ctx, campaignID, campaignUpdate, *currentDate)
	if err != nil {
		return nil, err
	}

	picURL, err := s.getPicURL(ctx, campaignID)
	if err == nil {
		// Set campaign pic url
		campaign.PicURL = &picURL
	}

	return campaign, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, advertiserID, campaignID uuid.UUID) error {
	// Check if advertiser exists
	_, err := s.advertiserRepo.GetByID(ctx, advertiserID)
	if err != nil {
		return domain.ErrAdvertiserNotFound
	}
	// Check if campaign exists
	_, err = s.repo.GetCampaignByID(ctx, campaignID)
	if err == pgx.ErrNoRows {
		return domain.ErrAdNotFound
	}
	err = s.repo.DeleteCampaign(ctx, campaignID)
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

func (s *CampaignService) getPicURL(ctx context.Context, campaignID uuid.UUID) (string, error) {
	// Get picture id from db
	picID, err := s.repo.GetCampaignPicID(ctx, campaignID)
	if err != nil || picID == "" {
		return "", nil
	}

	// Generate picture URL and return
	return s.fileRepo.GetFileLink(ctx, picID, s.minioPublicHost)
}

func validateTargeting(targeting domain.Targeting) bool {
	if targeting.AgeFrom != nil && *targeting.AgeFrom <= 0 {
		return false
	}
	if targeting.AgeTo != nil && *targeting.AgeTo >= 200 {
		return false
	}

	// Check if ageFrom is greater than ageTo
	if targeting.AgeFrom != nil && targeting.AgeTo != nil {
		if *targeting.AgeFrom > *targeting.AgeTo {
			return false
		}
	}

	if targeting.Gender != nil && !isValidGender(*targeting.Gender) {
		return false
	}

	return true
}

func isValidGender(gender string) bool {
	validGenders := map[string]struct{}{
		"MALE":   {},
		"FEMALE": {},
		"ALL":    {},
	}
	_, exists := validGenders[gender]
	return exists
}

func (s *CampaignService) validateModeration(ctx context.Context, adTitle, adText string) error {
	isModerated, err := s.checkModeration(ctx)
	if err != nil {
		return err
	}
	if isModerated {
		allText := fmt.Sprintf("Название: %s; Описание: %s", adTitle, adText)
		isAllowed, err := s.openAIService.ValidateAdText(ctx, allText)
		if err != nil {
			return err
		}
		if !isAllowed {
			return domain.ErrModerationNotPassed
		}
	}
	return nil
}
