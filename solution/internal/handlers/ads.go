package handlers

import (
	"encoding/json"
	"errors"
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

func (h *AdsHandler) GetAd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clientId, err := uuid.Parse(r.URL.Query().Get("client_id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID клиента")
		return
	}

	ad, err := h.service.GetAd(ctx, clientId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			WriteError(w, http.StatusNotFound, "Клиент не найден", "")
		case errors.Is(err, domain.ErrAdNotFound):
			WriteError(w, http.StatusNotFound, "Подходящая реклама не найдена", "")
		default:
			log.Printf("[INTERNAL ERROR] failed to get ad: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	json.NewEncoder(w).Encode(ad)
}

// Click godoc
//
//	@Summary	Фиксация перехода по рекламному объявлению
//	@Desciption	Фиксирует клик (переход) клиента по рекламному объявлению
//	@Tags		Ads
//	@Accept		json
//	@Param		click	body	domain.Click	true	"click"
//	@Param		adId	path	string			true	"UUID рекламного объявления (идентификатор кампании), по которому совершен клик"
//	@Success	204
//	@Failure	400	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/ads/{adId}/click [post]
func (h *AdsHandler) Click(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	adId, err := uuid.Parse(chi.URLParam(r, "adId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламы")
		return
	}

	var body domain.Click

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный JSON")
		return
	}

	err = h.service.Click(ctx, adId, body.ClientID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrAdNotFound):
			WriteError(w, http.StatusNotFound, err.Error(), "")
		default:
			log.Printf("[INTERNAL ERROR] failed to register click: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
