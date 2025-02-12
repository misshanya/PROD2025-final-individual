package domain

import "github.com/google/uuid"

type Advertiser struct {
	ID uuid.UUID `json:"advertiser_id"`
	Name string `json:"name"`
}