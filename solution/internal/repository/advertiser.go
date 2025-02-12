package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
)

type AdvertiserRepository struct {
	queries *storage.Queries
}

func NewAdvertiserRepository(queries *storage.Queries) *AdvertiserRepository {
	return &AdvertiserRepository{
		queries: queries,
	}
}

func (r *AdvertiserRepository) Create(ctx context.Context, advertisers []*domain.Advertiser) ([]*domain.Advertiser, error) {
	for _, advertiser := range advertisers {
		_, err := r.queries.GetAdvertiserByID(ctx, advertiser.ID)
		if err == nil {
			return []*domain.Advertiser{}, domain.ErrAdvertiserAlreadyExists
		}

		err = r.queries.CreateAdvertiser(ctx, storage.CreateAdvertiserParams{
			ID: advertiser.ID,
			Name: advertiser.Name,
		})
		if err != nil {
			return []*domain.Advertiser{}, err
		}
	}
	
	return advertisers, nil
}

func (r *AdvertiserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Advertiser, error) {
	advertiser, err := r.queries.GetAdvertiserByID(ctx, id)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrAdvertiserNotFound
	} else if err != nil {
		return nil, err
	}
	return &domain.Advertiser{
		ID: advertiser.ID,
		Name: advertiser.Name,
	}, nil
}

func (r *AdvertiserRepository) CreateMLScore(ctx context.Context, score *domain.MLScore) error {
	err := r.queries.CreateMLScore(ctx, storage.CreateMLScoreParams{
		ClientID: score.ClientID,
		AdvertiserID: score.AdvertiserID,
		Score: score.Score,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *AdvertiserRepository) GetMLScore(ctx context.Context, clientID, advertiserID uuid.UUID) (*domain.MLScore, error) {
	score, err := r.queries.GetMLScoreByIDs(ctx, storage.GetMLScoreByIDsParams{
		ClientID: clientID,
		AdvertiserID: advertiserID,
	})
	if err == nil {
		return &domain.MLScore{
			ClientID: score.ClientID,
			AdvertiserID: score.AdvertiserID,
			Score: score.Score,
		}, nil
	} else if err == pgx.ErrNoRows {
		return &domain.MLScore{}, pgx.ErrNoRows
	}
	return nil, err
}

func (r *AdvertiserRepository) UpdateMLScore(ctx context.Context, score *domain.MLScore) error {
	err := r.queries.UpdateMLScore(ctx, storage.UpdateMLScoreParams{
		Score: score.Score,
		ClientID: score.ClientID,
		AdvertiserID: score.AdvertiserID,
	})
	if err != nil {
		return err
	}
	return nil
}