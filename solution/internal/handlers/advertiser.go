package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
)

type AdvertiserHandler struct {
	service *app.AdvertiserService
}

func NewAdvertiserHandler(service *app.AdvertiserService) *AdvertiserHandler {
	return &AdvertiserHandler{
		service: service,
	}
}

func (h *AdvertiserHandler) CreateAdvertisers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var advertisers []*domain.Advertiser

	if err := json.NewDecoder(r.Body).Decode(&advertisers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newAdvertisers, err := h.service.CreateAdvertisers(ctx, advertisers)
	if err != nil {
		switch err {
		case domain.ErrAdvertiserAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
		case domain.ErrBadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("[INTERNAL ERROR] failed to create advertiser: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAdvertisers)
}

func (h *AdvertiserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		log.Printf("failed to convert advertiserID to uuid: %v", err)
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	advertiser, err := h.service.GetByID(ctx, advertiserID)
	if err != nil {
		switch err {
		case domain.ErrAdvertiserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			log.Printf("[INTERNAL ERROR] failed to get advertiser: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		}
	}

	json.NewEncoder(w).Encode(advertiser)
}