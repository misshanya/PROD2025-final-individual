// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package storage

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Advertiser struct {
	ID   uuid.UUID
	Name string
}

type Campaign struct {
	ID                uuid.UUID
	AdvertiserID      uuid.UUID
	ImpressionsLimit  int64
	ClicksLimit       int64
	CostPerImpression pgtype.Numeric
	CostPerClick      pgtype.Numeric
	AdTitle           string
	AdText            string
	StartDate         int32
	EndDate           int32
}

type CampaignsTargeting struct {
	ID         uuid.UUID
	CampaignID uuid.UUID
	Gender     pgtype.Text
	AgeFrom    pgtype.Int4
	AgeTo      pgtype.Int4
	Location   pgtype.Text
}

type MlScore struct {
	ClientID     uuid.UUID
	AdvertiserID uuid.UUID
	Score        int32
}

type User struct {
	ID       uuid.UUID
	Login    string
	Age      int32
	Location string
	Gender   string
}
