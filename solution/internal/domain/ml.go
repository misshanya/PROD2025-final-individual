package domain

import (
	"context"

	"github.com/google/uuid"
)

type MLService interface {
	ValidateAdText(ctx context.Context, text string) (bool, error)
	GenerateAdText(ctx context.Context, advertiserName, adTitle string) (string, error)
}

type MLScore struct {
	ClientID     uuid.UUID `json:"client_id"`
	AdvertiserID uuid.UUID `json:"advertiser_id"`
	Score        int32     `json:"score"`
}
