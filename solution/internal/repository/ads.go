package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *AdsRepository) Impression(ctx context.Context, adId, clientId uuid.UUID) error {
	_, err := r.queries.CreateImpression(ctx, storage.CreateImpressionParams{
		CampaignID: adId,
		ClientID: clientId,
	})
	return err
}

func (r *AdsRepository) Click(ctx context.Context, adId, clientId uuid.UUID) error {
	isClicked, err := r.queries.IsClicked(ctx, storage.IsClickedParams{
		CampaignID: adId,
		ClientID: clientId,
	})
	if isClicked == 1 {
		return nil
	} else if err != nil && err != pgx.ErrNoRows {
		return err
	}
	_, err = r.queries.CreateClick(ctx, storage.CreateClickParams{
		CampaignID: adId,
		ClientID: clientId,
	})
	return err
}