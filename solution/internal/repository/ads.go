package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
)

type AdsRepository struct {
	queries *storage.Queries
}

func NewAdsRepository(queries *storage.Queries) *AdsRepository {
	return &AdsRepository{
		queries: queries,
	}
}

func (r *AdsRepository) GetRelativeAd(ctx context.Context, clientId uuid.UUID, currentDate int32) (*domain.UserAd, error) {
	client, err := r.queries.GetUserByID(ctx, clientId)
	if err != nil {
		return nil, err
	}
	ad, err := r.queries.GetRelativeAd(ctx, storage.GetRelativeAdParams{
		ClientID: clientId,
		Gender:   client.Gender,
		Age:      client.Age,
		Location: client.Location,
		CurDate:  currentDate,
	})
	if err != nil {
		return nil, err
	}
	return &domain.UserAd{
		AdId:         ad.ID,
		AdTitle:      ad.AdTitle,
		AdText:       ad.AdText,
		AdvertiserID: ad.AdvertiserID,
	}, nil
}

func (r *AdsRepository) Impression(ctx context.Context, adId, clientId uuid.UUID) error {
	_, err := r.queries.CreateImpression(ctx, storage.CreateImpressionParams{
		CampaignID: adId,
		ClientID:   clientId,
	})
	return err
}

func (r *AdsRepository) Click(ctx context.Context, adId, clientId uuid.UUID) error {
	isClicked, err := r.queries.IsClicked(ctx, storage.IsClickedParams{
		CampaignID: adId,
		ClientID:   clientId,
	})
	if isClicked == 1 {
		return nil
	} else if err != nil && err != pgx.ErrNoRows {
		return err
	}
	_, err = r.queries.CreateClick(ctx, storage.CreateClickParams{
		CampaignID: adId,
		ClientID:   clientId,
	})
	return err
}
