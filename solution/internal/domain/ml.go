package domain

import "github.com/google/uuid"

type MLScore struct {
	UserID uuid.UUID `json:"user_id"`
	AdvertiserID uuid.UUID `json:"advertiser_id"`
}