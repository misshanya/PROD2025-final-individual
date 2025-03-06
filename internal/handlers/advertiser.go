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

type AdvertiserHandler struct {
	service *app.AdvertiserService
}

func NewAdvertiserHandler(service *app.AdvertiserService) *AdvertiserHandler {
	return &AdvertiserHandler{
		service: service,
	}
}

// CreateAdvertisers godoc
//
//	@Summary		Создание/обновление рекламодателей
//	@Description	Создает или обновляет рекламодателей
//	@Tags			Advertisers
//	@Accept			json
//	@Param			CreateAdvertisers	body	[]domain.Advertiser	true	"CampaignRequest"
//	@Produce		json
//	@Success		201	{object}	[]domain.Advertiser
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/bulk [post]
func (h *AdvertiserHandler) CreateAdvertisers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var advertisers []*domain.Advertiser

	if err := json.NewDecoder(r.Body).Decode(&advertisers); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	newAdvertisers, err := h.service.CreateUpdateAdvertisers(ctx, advertisers)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to create advertiser: %v", err)
		WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAdvertisers)
}

// GetByID godoc
//
//	@Summary		Получение рекламодателя по ID
//	@Description	Возвращает информацию о рекламодателе по его ID
//	@Tags			Advertisers
//	@Produce		json
//	@Success		200	{object}	domain.Advertiser
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId} [get]
func (h *AdvertiserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	advertiser, err := h.service.GetByID(ctx, advertiserID)
	if err != nil {
		switch err {
		case domain.ErrAdvertiserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			log.Printf("[INTERNAL ERROR] failed to get advertiser: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	json.NewEncoder(w).Encode(advertiser)
}

// CreateUpdateMLScore godoc
//
//	@Summary		Добавление или обновление ML скора
//	@Description	Добавляет или обновляет ML скор для указанной пары клиент-рекламодатель
//	@Tags			Advertisers
//	@Accept			json
//	@Param			MLScore	body	domain.MLScore	true	"MLScore"
//	@Produce		json
//	@Success		200	{object}	domain.MLScore
//	@Failure		404	{object}	ErrorResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/ml-scores [post]
func (h *AdvertiserHandler) CreateUpdateMLScore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var score *domain.MLScore

	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	newScore, err := h.service.CreateUpdateMLScore(ctx, score)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdvertiserNotFound):
			WriteError(w, http.StatusNotFound, "Рекламодатель не найден", "")
		case errors.Is(err, domain.ErrUserNotFound):
			WriteError(w, http.StatusNotFound, "Клиент не найден", "")
		default:
			log.Printf("failed to create or update ml score: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}

		return
	}

	json.NewEncoder(w).Encode(newScore)
}
