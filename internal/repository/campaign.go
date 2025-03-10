package repository

import (
	"context"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
)

type CampaignRepository struct {
	queries *storage.Queries
	dbConn  *pgxpool.Pool
}

func NewCampaignRepository(queries *storage.Queries, dbConn *pgxpool.Pool) *CampaignRepository {
	return &CampaignRepository{
		queries: queries,
		dbConn:  dbConn,
	}
}

func (r *CampaignRepository) CreateCampaign(ctx context.Context, advertiserID uuid.UUID, campaignRequest *domain.CampaignRequest, currentDate int) (*domain.Campaign, error) {
	// Convert cost per impression and cost per click to pgtype.Numeric
	costPerImpression, err := convertCostToNumeric(campaignRequest.CostPerImpression)
	if err != nil {
		return nil, err
	}
	costPerClick, err := convertCostToNumeric(campaignRequest.CostPerClick)
	if err != nil {
		return nil, err
	}

	if campaignRequest.StartDate < int32(currentDate) || campaignRequest.EndDate < int32(currentDate) || campaignRequest.EndDate < campaignRequest.StartDate {
		return nil, domain.ErrBadRequest
	}

	// Create campaign and targeting in transaction
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	campaignDB, err := qtx.CreateCampaign(ctx, storage.CreateCampaignParams{
		AdvertiserID:      advertiserID,
		ImpressionsLimit:  campaignRequest.ImpressionsLimit,
		ClicksLimit:       campaignRequest.ClicksLimit,
		CostPerImpression: costPerImpression,
		CostPerClick:      costPerClick,
		AdTitle:           campaignRequest.AdTitle,
		AdText:            campaignRequest.AdText,
		StartDate:         campaignRequest.StartDate,
		EndDate:           campaignRequest.EndDate,
	})
	if err != nil {
		return nil, err
	}

	createCampaignTargetingParams := buildCampaignTargetingParams(campaignRequest.Targeting)
	createCampaignTargetingParams.CampaignID = campaignDB.ID

	targetingDB, err := qtx.CreateCampaignTargeting(ctx, createCampaignTargetingParams)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)

	campaign := domain.Campaign{
		ID:                campaignDB.ID,
		AdvertiserID:      campaignDB.AdvertiserID,
		ImpressionsLimit:  campaignDB.ImpressionsLimit,
		ClicksLimit:       campaignDB.ClicksLimit,
		CostPerImpression: campaignRequest.CostPerImpression,
		CostPerClick:      campaignRequest.CostPerClick,
		AdTitle:           campaignDB.AdTitle,
		AdText:            campaignDB.AdText,
		StartDate:         campaignDB.StartDate,
		EndDate:           campaignDB.EndDate,
		Targeting:         convertDBTargetingToDomain(targetingDB),
	}
	return &campaign, nil
}

func (r *CampaignRepository) SetCampaignPicture(ctx context.Context, campaignID uuid.UUID, picID string) error {
	return r.queries.SetCampaignPicture(ctx, storage.SetCampaignPictureParams{
		PictureID:  picID,
		CampaignID: campaignID,
	})
}

func (r *CampaignRepository) GetCampaignsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID, size, offset int) ([]domain.Campaign, error) {
	campaignsDB, err := r.queries.GetCampaignsWithTargetingByAdvertiserID(ctx, storage.GetCampaignsWithTargetingByAdvertiserIDParams{
		Limit:        int32(size),
		Offset:       int32(offset),
		AdvertiserID: advertiserID,
	})
	if err != nil {
		return nil, err
	}

	campaigns := make([]domain.Campaign, len(campaignsDB))

	for i, campaignDB := range campaignsDB {
		// Convert cost per impression and cost per click to float64 values
		var costPerImpression float64
		costPerImpressionFloatDB, err := campaignDB.CostPerImpression.Float64Value()
		if err != nil {
			return nil, err
		}
		costPerImpression = costPerImpressionFloatDB.Float64

		var costPerClick float64
		costPerClickFloatDB, err := campaignDB.CostPerClick.Float64Value()
		if err != nil {
			return nil, err
		}
		costPerClick = costPerClickFloatDB.Float64

		targeting := domain.Targeting{}

		if campaignDB.Gender.Valid {
			targeting.Gender = &campaignDB.Gender.String
		}
		if campaignDB.AgeFrom.Valid {
			targeting.AgeFrom = &campaignDB.AgeFrom.Int32
		}
		if campaignDB.AgeTo.Valid {
			targeting.AgeTo = &campaignDB.AgeTo.Int32
		}
		if campaignDB.Location.Valid {
			targeting.Location = &campaignDB.Location.String
		}
		campaigns[i] = domain.Campaign{
			ID:                campaignDB.ID,
			AdvertiserID:      advertiserID,
			ImpressionsLimit:  campaignDB.ImpressionsLimit,
			ClicksLimit:       campaignDB.ClicksLimit,
			CostPerImpression: costPerImpression,
			CostPerClick:      costPerClick,
			AdTitle:           campaignDB.AdTitle,
			AdText:            campaignDB.AdText,
			StartDate:         campaignDB.StartDate,
			EndDate:           campaignDB.EndDate,
			Targeting:         targeting,
		}
	}

	return campaigns, nil
}

func (r *CampaignRepository) GetCampaignPicID(ctx context.Context, campaignID uuid.UUID) (string, error) {
	picIDPg, err := r.queries.GetCampaignPicID(ctx, campaignID)
	if err != nil || !picIDPg.Valid {
		return "", err
	}
	return picIDPg.String, nil
}

func (r *CampaignRepository) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*domain.Campaign, error) {
	campaignDB, err := r.queries.GetCampaignWithTargetingByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Convert cost per impression and cost per click to float64 values
	var costPerImpression float64
	costPerImpressionFloatDB, err := campaignDB.CostPerImpression.Float64Value()
	if err != nil {
		return nil, err
	}
	costPerImpression = costPerImpressionFloatDB.Float64

	var costPerClick float64
	costPerClickFloatDB, err := campaignDB.CostPerClick.Float64Value()
	if err != nil {
		return nil, err
	}
	costPerClick = costPerClickFloatDB.Float64

	targeting := domain.Targeting{}

	if campaignDB.Gender.Valid {
		targeting.Gender = &campaignDB.Gender.String
	}
	if campaignDB.AgeFrom.Valid {
		targeting.AgeFrom = &campaignDB.AgeFrom.Int32
	}
	if campaignDB.AgeTo.Valid {
		targeting.AgeTo = &campaignDB.AgeTo.Int32
	}
	if campaignDB.Location.Valid {
		targeting.Location = &campaignDB.Location.String
	}

	return &domain.Campaign{
		ID:                campaignDB.ID,
		AdvertiserID:      campaignDB.AdvertiserID,
		ImpressionsLimit:  campaignDB.ImpressionsLimit,
		ClicksLimit:       campaignDB.ClicksLimit,
		CostPerImpression: costPerImpression,
		CostPerClick:      costPerClick,
		AdTitle:           campaignDB.AdTitle,
		AdText:            campaignDB.AdText,
		StartDate:         campaignDB.StartDate,
		EndDate:           campaignDB.EndDate,
		Targeting:         targeting,
	}, nil
}

func (r *CampaignRepository) UpdateCampaign(ctx context.Context, campaignID uuid.UUID, campaignUpdate domain.CampaignUpdateRequest, currentDate int) (*domain.Campaign, error) {
	// Convert cost per impression and cost per click to pgtype.Numeric
	costPerImpression, err := convertCostToNumeric(campaignUpdate.CostPerImpression)
	if err != nil {
		return nil, err
	}
	costPerClick, err := convertCostToNumeric(campaignUpdate.CostPerClick)
	if err != nil {
		return nil, err
	}

	var started bool
	existingCampaignDB, err := r.queries.GetCampaignWithTargetingByID(ctx, campaignID)
	started = existingCampaignDB.StartDate <= int32(currentDate)
	log.Printf("updating campaign %s, started: %v", existingCampaignDB.ID, started)

	var impressionsLimit, clicksLimit int64
	var startDate, endDate int32
	if started {
		impressionsLimit = existingCampaignDB.ImpressionsLimit
		clicksLimit = existingCampaignDB.ClicksLimit
		startDate = existingCampaignDB.StartDate
		endDate = existingCampaignDB.EndDate
	} else {
		impressionsLimit = campaignUpdate.ImpressionsLimit
		clicksLimit = campaignUpdate.ClicksLimit
		startDate = campaignUpdate.StartDate
		endDate = campaignUpdate.EndDate
	}

	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	campaignDB, err := qtx.UpdateCampaign(ctx, storage.UpdateCampaignParams{
		CampaignID:        campaignID,
		ImpressionsLimit:  impressionsLimit,
		ClicksLimit:       clicksLimit,
		CostPerImpression: costPerImpression,
		CostPerClick:      costPerClick,
		AdTitle:           campaignUpdate.AdTitle,
		AdText:            campaignUpdate.AdText,
		StartDate:         startDate,
		EndDate:           endDate,
	})
	if err != nil {
		return nil, err
	}

	updateCampaignTargetingParams := buildUpdateCampaignTargetingParams(campaignUpdate.Targeting)
	updateCampaignTargetingParams.CampaignID = campaignDB.ID

	targetingDB, err := qtx.UpdateCampaignTargeting(ctx, updateCampaignTargetingParams)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)

	campaign := domain.Campaign{
		ID:                campaignDB.ID,
		AdvertiserID:      campaignDB.AdvertiserID,
		ImpressionsLimit:  campaignDB.ImpressionsLimit,
		ClicksLimit:       campaignDB.ClicksLimit,
		CostPerImpression: campaignUpdate.CostPerImpression,
		CostPerClick:      campaignUpdate.CostPerClick,
		AdTitle:           campaignDB.AdTitle,
		AdText:            campaignDB.AdText,
		StartDate:         campaignDB.StartDate,
		EndDate:           campaignDB.EndDate,
		Targeting:         convertDBTargetingToDomain(targetingDB),
	}
	return &campaign, nil
}

func (r *CampaignRepository) DeleteCampaign(ctx context.Context, campaignID uuid.UUID) error {
	err := r.queries.DeleteCampaignByID(ctx, campaignID)
	return err
}

func convertCostToNumeric(cost float64) (pgtype.Numeric, error) {
	var num pgtype.Numeric
	costStr := strconv.FormatFloat(cost, 'f', -1, 64)
	if err := num.Scan(costStr); err != nil {
		return num, err
	}
	return num, nil
}

func buildCampaignTargetingParams(targeting domain.Targeting) storage.CreateCampaignTargetingParams {
	params := storage.CreateCampaignTargetingParams{}
	if targeting.Gender != nil {
		params.Gender = pgtype.Text{String: *targeting.Gender, Valid: true}
	}
	if targeting.AgeFrom != nil {
		params.AgeFrom = pgtype.Int4{Int32: *targeting.AgeFrom, Valid: true}
	}
	if targeting.AgeTo != nil {
		params.AgeTo = pgtype.Int4{Int32: *targeting.AgeTo, Valid: true}
	}
	if targeting.Location != nil {
		params.Location = pgtype.Text{String: *targeting.Location, Valid: true}
	}
	return params
}

func buildUpdateCampaignTargetingParams(targeting domain.Targeting) storage.UpdateCampaignTargetingParams {
	params := storage.UpdateCampaignTargetingParams{}
	if targeting.Gender != nil {
		params.Gender = pgtype.Text{String: *targeting.Gender, Valid: true}
	}
	if targeting.AgeFrom != nil {
		params.AgeFrom = pgtype.Int4{Int32: *targeting.AgeFrom, Valid: true}
	}
	if targeting.AgeTo != nil {
		params.AgeTo = pgtype.Int4{Int32: *targeting.AgeTo, Valid: true}
	}
	if targeting.Location != nil {
		params.Location = pgtype.Text{String: *targeting.Location, Valid: true}
	}
	return params
}

func convertDBTargetingToDomain(targetingDB storage.CampaignsTargeting) domain.Targeting {
	targeting := domain.Targeting{}
	if targetingDB.Gender.Valid {
		targeting.Gender = &targetingDB.Gender.String
	}
	if targetingDB.AgeFrom.Valid {
		targeting.AgeFrom = &targetingDB.AgeFrom.Int32
	}
	if targetingDB.AgeTo.Valid {
		targeting.AgeTo = &targetingDB.AgeTo.Int32
	}
	if targetingDB.Location.Valid {
		targeting.Location = &targetingDB.Location.String
	}
	return targeting
}
