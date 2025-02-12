package domain

import "github.com/google/uuid"

type Advertiser struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
}