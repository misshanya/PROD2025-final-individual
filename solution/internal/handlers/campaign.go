package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
)

type CampaignHandler struct {
	service *app.CampaignService
}

func NewCampaignHandler(service *app.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		service: service,
	}
}

func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	var campaignRequest *domain.CampaignRequest

	if err := json.NewDecoder(r.Body).Decode(&campaignRequest); err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	campaign, err := h.service.CreateCampaign(ctx, advertiserID, campaignRequest)
	if err != nil {
		switch err {
		case domain.ErrBadRequest, domain.ErrModerationNotPassed:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("[INTERNAL ERROR] failed to create campaign: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(campaign)
}

func (h *CampaignHandler) GetCampaignsByAdvertiserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	var size, page int
	sizeStr := r.URL.Query().Get("size")
	if sizeStr == "" {
		size = 10
	} else {
		sizeTmp, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}
		size = sizeTmp
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = 0
	} else {
		pageTmp, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}
		page = pageTmp
	}

	campaigns, err := h.service.GetCampaignsByAdvertiserID(ctx, advertiserID, size, page)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to get campaigns: %v", err)
		http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(campaigns)
}

func (h *CampaignHandler) GetCampaignByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	campaign, err := h.service.GetCampaignByID(ctx, campaignID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			log.Printf("[INTERNAL ERROR] failed to get campaign by id: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(campaign)
}

func (h *CampaignHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	var campaignUpdate domain.CampaignUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&campaignUpdate); err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	newCampaign, err := h.service.UpdateCampaign(ctx, campaignID, campaignUpdate)
	if err != nil {
		switch err {
		case domain.ErrModerationNotPassed:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("[INTERNAL ERROR] failed to update campaign: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(newCampaign)
}

func (h *CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCampaign(ctx, campaignID)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to delete campaign: %v", err)
		http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CampaignHandler) GenerateAdText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var GenerateAdTextRequest *domain.GenerateAdTextRequest

	if err := json.NewDecoder(r.Body).Decode(&GenerateAdTextRequest); err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	adText, err := h.service.GenerateAdText(ctx, GenerateAdTextRequest.AdvertiserName, GenerateAdTextRequest.AdTitle)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to generate ad text: %v", err)
		http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	response := domain.GenerateAdTextResponse{
		AdText: adText,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *CampaignHandler) SwitchModeration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	newIsModerated, err := h.service.SwitchModeration(ctx)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to switch moderation: %v", err)
		http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	response := domain.SwitchModerationResponse{
		IsModerated: newIsModerated,
	}

	json.NewEncoder(w).Encode(response)
}
