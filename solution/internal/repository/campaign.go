package repository

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
)


type CampaignRepository struct {
	queries *storage.Queries
	dbConn *pgx.Conn
}

func NewCampaignRepository(queries *storage.Queries, dbConn *pgx.Conn) *CampaignRepository {
	return &CampaignRepository{
		queries: queries,
		dbConn: dbConn,
	}
}

func (r *CampaignRepository) CreateCampaign(ctx context.Context, advertiserID uuid.UUID, campaignRequest *domain.CampaignRequest) (*domain.Campaign, error) {
	// Convert cost per impression and cost per click to pgtype.Numeric
	var costPerImpression pgtype.Numeric
	var costPerClick pgtype.Numeric
	costPerImpression.Scan(strconv.FormatFloat(campaignRequest.CostPerImpression, 'f', -1, 64))
	costPerClick.Scan(strconv.FormatFloat(campaignRequest.CostPerClick, 'f', -1, 64))

	// Create campaign and targeting in transaction
	tx, err := r.dbConn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	campaignDB, err := qtx.CreateCampaign(ctx, storage.CreateCampaignParams{
		AdvertiserID: advertiserID,
		ImpressionsLimit: campaignRequest.ImpressionsLimit,
		ClicksLimit: campaignRequest.ClicksLimit,
		CostPerImpression: costPerImpression,
		CostPerClick: costPerClick,
		AdTitle: campaignRequest.AdTitle,
		AdText: campaignRequest.AdText,
		StartDate: campaignRequest.StartDate,
		EndDate: campaignRequest.EndDate,
	})
	if err != nil {
		return nil, err
	}

	createCampaignTargetingParams := storage.CreateCampaignTargetingParams{
		CampaignID: campaignDB.ID,
	}

	if campaignRequest.Targeting.Gender != nil {
		createCampaignTargetingParams.Gender = pgtype.Text{String: *campaignRequest.Targeting.Gender, Valid: true}
	}
	if campaignRequest.Targeting.AgeFrom != nil {
		createCampaignTargetingParams.AgeFrom = pgtype.Int4{Int32: *campaignRequest.Targeting.AgeFrom, Valid: true}
	}
	if campaignRequest.Targeting.AgeTo != nil {
		createCampaignTargetingParams.AgeTo = pgtype.Int4{Int32: *campaignRequest.Targeting.AgeTo, Valid: true}
	}
	if campaignRequest.Targeting.Location != nil {
		createCampaignTargetingParams.Location = pgtype.Text{String: *campaignRequest.Targeting.Location, Valid: true}
	}

	targetingDB, err := qtx.CreateCampaignTargeting(ctx, createCampaignTargetingParams)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)

	// Convert db responses to my structures
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
	campaign := domain.Campaign{
		ID: campaignDB.ID,
		AdvertiserID: campaignDB.AdvertiserID,
		ImpressionsLimit: campaignDB.ImpressionsLimit,
		ClicksLimit: campaignDB.ClicksLimit,
		CostPerImpression: campaignRequest.CostPerImpression,
		CostPerClick: campaignRequest.CostPerClick,
		AdTitle: campaignDB.AdTitle,
		AdText: campaignDB.AdText,
		StartDate: campaignDB.StartDate,
		EndDate: campaignDB.EndDate,
		Targeting: targeting,
	}
	return &campaign, nil
}

func (r *CampaignRepository) GetCampaignsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID, size, offset int) ([]domain.Campaign, error) {
	campaignsDB, err := r.queries.GetCampaignsWithTargetingByAdvertiserID(ctx, storage.GetCampaignsWithTargetingByAdvertiserIDParams{
		Limit: int32(size),
		Offset: int32(offset),
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
			ID: campaignDB.ID,
			AdvertiserID: advertiserID,
			ImpressionsLimit: campaignDB.ImpressionsLimit,
			ClicksLimit: campaignDB.ClicksLimit,
			CostPerImpression: costPerImpression,
			CostPerClick: costPerClick,
			AdTitle: campaignDB.AdTitle,
			AdText: campaignDB.AdText,
			StartDate: campaignDB.StartDate,
			EndDate: campaignDB.EndDate,
			Targeting: targeting,
		}
	}

	return campaigns, nil
}

func (r *CampaignRepository) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*domain.Campaign, error) {
	campaignDB, err := r.queries.GetCampaignWithTargetingByID(ctx, campaignID)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	} else if err != nil {
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
		ID: campaignDB.ID,
		AdvertiserID: campaignDB.AdvertiserID,
		ImpressionsLimit: campaignDB.ImpressionsLimit,
		ClicksLimit: campaignDB.ClicksLimit,
		CostPerImpression: costPerImpression,
		CostPerClick: costPerClick,
		AdTitle: campaignDB.AdTitle,
		AdText: campaignDB.AdText,
		StartDate: campaignDB.StartDate,
		EndDate: campaignDB.EndDate,
		Targeting: targeting,
	}, nil
}