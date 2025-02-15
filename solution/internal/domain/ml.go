package domain

import "github.com/google/uuid"

type MLScore struct {
	ClientID     uuid.UUID `json:"client_id"`
	AdvertiserID uuid.UUID `json:"advertiser_id"`
	Score        int32     `json:"score"`
}
