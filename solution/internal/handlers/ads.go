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

type AdsHandler struct {
	service *app.AdsService
}

func NewAdsHandler(service *app.AdsService) *AdsHandler {
	return &AdsHandler{
		service: service,
	}
}

func (h *AdsHandler) Click(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	adId, err := uuid.Parse(chi.URLParam(r, "adId"))
	if err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	var body struct {
		ClientID uuid.UUID `json:"client_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Click(ctx, adId, body.ClientID)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to register click: %v", err)
		http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}