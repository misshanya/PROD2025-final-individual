package domain

import "github.com/google/uuid"

type Campaign struct {
	ID                uuid.UUID `json:"campaign_id"`
	AdvertiserID      uuid.UUID `json:"advertiser_id"`
	ImpressionsLimit  int64     `json:"impressions_limit"`
	ClicksLimit       int64     `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int32     `json:"start_date"`
	EndDate           int32     `json:"end_date"`
	Targeting         Targeting `json:"targeting"`
}

type CampaignRequest struct {
	ImpressionsLimit  int64     `json:"impressions_limit"`
	ClicksLimit       int64     `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int32     `json:"start_date"`
	EndDate           int32     `json:"end_date"`
	Targeting         Targeting `json:"targeting"`
}

type Targeting struct {
	Gender   *string `json:"gender,omitempty"`
	AgeFrom  *int32  `json:"age_from,omitempty"`
	AgeTo    *int32  `json:"age_to,omitempty"`
	Location *string `json:"location,omitempty"`
}

type CampaignUpdateRequest struct {
	ImpressionsLimit  int64     `json:"impressions_limit"`
	ClicksLimit       int64     `json:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click"`
	AdTitle           string    `json:"ad_title"`
	AdText            string    `json:"ad_text"`
	StartDate         int32     `json:"start_date"`
	EndDate           int32     `json:"end_date"`
	Targeting         Targeting `json:"targeting"`
}

type GenerateAdTextRequest struct {
	AdvertiserName string `json:"advertiser_name"`
	AdTitle        string `json:"ad_title"`
}

type GenerateAdTextResponse struct {
	AdText string `json:"ad_text"`
}

type SwitchModerationResponse struct {
	IsModerated bool `json:"is_moderated"`
}
