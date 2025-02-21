package domain

import "github.com/google/uuid"

type Click struct {
	ClientID uuid.UUID `json:"client_id"`
}

type UserAd struct {
	AdId         uuid.UUID `json:"ad_id"`
	AdTitle      string    `json:"ad_title"`
	AdText       string    `json:"ad_text"`
	AdvertiserID uuid.UUID `json:"advertiser_id"`
}
