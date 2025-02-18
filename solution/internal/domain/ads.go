package domain

import "github.com/google/uuid"

type Click struct {
	ClientID uuid.UUID `json:"client_id"`
}
