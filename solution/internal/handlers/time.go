package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/app"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
)

type TimeHandler struct {
	service *app.TimeService
}

func NewTimeHandler(service *app.TimeService) *TimeHandler {
	return &TimeHandler{
		service: service,
	}
}

func (h *TimeHandler) SetCurrentDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		CurrentDate int `json:"current_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SetCurrentDate(ctx, body.CurrentDate); err != nil {
		switch err {
		case domain.ErrNewDateLowerThanCurrent:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			log.Printf("[INTERNAL ERROR] failed to set current date: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(body)
}
