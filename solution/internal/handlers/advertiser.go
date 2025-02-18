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

// CreateAdvertisers godoc
//
//	@Summary		Создание/обновление рекламодателей
//	@Description	Создает или обновляет рекламодателей
//	@Tags			Advertisers
//	@Accept			json
//	@Param			CreateAdvertisers	body	[]domain.Advertiser	true	"CampaignRequest"
//	@Produce		json
//	@Success		201	{object}	[]domain.Advertiser
//	@Failure		400	{string}	BadRequest	"Bad request"
//	@Failure		500	{string}	string		"Internal Server Error"
//	@Router			/advertisers/bulk [post]
func (h *AdvertiserHandler) CreateAdvertisers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var advertisers []*domain.Advertiser

	if err := json.NewDecoder(r.Body).Decode(&advertisers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newAdvertisers, err := h.service.CreateUpdateAdvertisers(ctx, advertisers)
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

// GetByID godoc
//
//	@Summary		Получение рекламодателя по ID
//	@Description	Возвращает информацию о рекламодателе по его ID
//	@Tags			Advertisers
//	@Produce		json
//	@Success		200	{object}	domain.Advertiser
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/advertisers/{advertiserId} [get]
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

// CreateUpdateMLScore godoc
//
//	@Summary		Добавление или обновление ML скора
//	@Description	Добавляет или обновляет ML скор для указанной пары клиент-рекламодатель
//	@Tags			Advertisers
//	@Accept			json
//	@Param			MLScore	body	domain.MLScore	true	"MLScore"
//	@Produce		json
//	@Success		200
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/ml-scores [post]
func (h *AdvertiserHandler) CreateUpdateMLScore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var score *domain.MLScore

	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newScore, err := h.service.CreateUpdateMLScore(ctx, score)
	if err != nil {
		log.Printf("failed to create or update ml score: %v", err)
		http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newScore)
}
