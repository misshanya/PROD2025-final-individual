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