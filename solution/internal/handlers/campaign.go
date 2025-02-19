package handlers

import (
	"encoding/json"
	"errors"
	"io"
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

// CreateCampaign godoc
//
//	@Summary		Создание кампании
//	@Description	Создает рекламную кампанию
//	@Tags			Campaigns
//	@Accept			json
//	@Param			CreateCampaign	body	domain.CampaignRequest	true	"CampaignRequest"
//	@Param			advertiserId	path	string					true	"Advertiser ID"
//	@Produce		json
//	@Success		200	{object}	domain.Campaign
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId}/campaigns [post]
func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	var campaignRequest *domain.CampaignRequest

	if err := json.NewDecoder(r.Body).Decode(&campaignRequest); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	campaign, err := h.service.CreateCampaign(ctx, advertiserID, campaignRequest)
	if err != nil {
		switch err {
		case domain.ErrBadRequest, domain.ErrModerationNotPassed:
			WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		default:
			log.Printf("[INTERNAL ERROR] failed to create campaign: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(campaign)
}

// SetCampaignPicture godoc
//
//	@Summary		Добавление картинки к рекламной кампании
//	@Description	Добавляет/обновляет изображение рекламной кампании
//	@Tags			Campaigns
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			advertiserId	path		string	true	"UUID рекламодателя"
//	@Param			campaignId		path		string	true	"UUID рекламной кампании"
//	@Param			uploadfile		formData	file	true	"Файл изображения для загрузки"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId}/campaigns/{campaignId}/picture [post]
func (h *CampaignHandler) SetCampaignPicture(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламной кампании")
		return
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Failed to retrieve file: %v", err)
		WriteError(w, http.StatusBadRequest, "Ошибка получения файла", "")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Ошибка получения файла", "некорректный файл")
		return
	}

	err = h.service.SetCampaignPicture(ctx, advertiserID, campaignID, handler.Filename, fileBytes)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdvertiserNotFound):
			WriteError(w, http.StatusNotFound, "Рекламодатель не найден", "")
		case errors.Is(err, domain.ErrAdNotFound):
			WriteError(w, http.StatusNotFound, "Рекламная кампания не найдена", "")
		default:
			log.Printf("[INTERNAL ERROR] failed to set campaign picture: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}
}

// GetCampaignsByAdvertiserID godoc
//
//	@Summary		Получение кампаний рекламодателя
//	@Description	Возвращает кампании рекламодателя по его ID
//	@Tags			Campaigns
//	@Produce		json
//	@Param			advertiserId	path		string	true	"ID рекламодателя"
//	@Success		200				{object}	[]domain.Campaign
//	@Failure		400				{object}	ErrorResponse
//	@Failure		404				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId}/campaigns [get]
func (h *CampaignHandler) GetCampaignsByAdvertiserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	var size, page int
	sizeStr := r.URL.Query().Get("size")
	if sizeStr == "" {
		size = 10
	} else {
		sizeTmp, err := strconv.Atoi(sizeStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный size")
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
			WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный page")
			return
		}
		page = pageTmp
	}

	campaigns, err := h.service.GetCampaignsByAdvertiserID(ctx, advertiserID, size, page)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdvertiserNotFound):
			WriteError(w, http.StatusNotFound, "Рекламодатель не найден", "")
		default:
			log.Printf("[INTERNAL ERROR] failed to get campaigns: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	json.NewEncoder(w).Encode(campaigns)
}

// GetCampaignByID godoc
//
//	@Summary		Получение кампании
//	@Description	Возвращает кампанию по ее ID
//	@Tags			Campaigns
//	@Produce		json
//	@Param			advertiserId	path		string	true	"ID рекламодателя"
//	@Param			campaignId		path		string	true	"ID рекламной кампании"
//	@Success		200				{object}	domain.Campaign
//	@Failure		400				{object}	ErrorResponse
//	@Failure		404				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId}/campaigns/{campaignId} [get]
func (h *CampaignHandler) GetCampaignByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламной кампании")
		return
	}

	campaign, err := h.service.GetCampaignByID(ctx, advertiserID, campaignID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdvertiserNotFound):
			WriteError(w, http.StatusNotFound, "Рекламодатель не найден", "")
		case errors.Is(err, domain.ErrAdNotFound):
			WriteError(w, http.StatusNotFound, "Рекламная компания не найдена", "")
		default:
			log.Printf("[INTERNAL ERROR] failed to get campaign by id: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	json.NewEncoder(w).Encode(campaign)
}

// UpdateCampaign godoc
//
//	@Summary		Обновление кампании
//	@Description	Обновляет разрешённые параметры рекламной кампании
//	@Tags			Campaigns
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domain.Campaign
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId}/campaigns/{campaignId} [put]
func (h *CampaignHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламной кампании")
		return
	}

	var campaignUpdate domain.CampaignUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&campaignUpdate); err != nil {
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	newCampaign, err := h.service.UpdateCampaign(ctx, advertiserID, campaignID, campaignUpdate)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdvertiserNotFound):
			WriteError(w, http.StatusNotFound, "Рекламодатель не найден", "")
		case errors.Is(err, domain.ErrAdNotFound):
			WriteError(w, http.StatusNotFound, "Рекламная кампания не найдена", "")
		case errors.Is(err, domain.ErrModerationNotPassed):
			WriteError(w, http.StatusBadRequest, "Некорректный запрос", "модерация не пройдена")
		default:
			log.Printf("[INTERNAL ERROR] failed to update campaign: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	json.NewEncoder(w).Encode(newCampaign)
}

// DeleteCampaign godoc
//
//	@Summary		Удаление рекламной кампании
//	@Description	Удаляет рекламную кампанию по ее ID
//	@Tags			Campaigns
//	@Success		204
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/{advertiserId}/campaigns/{campaignId} [delete]
func (h *CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advertiserID, err := uuid.Parse(chi.URLParam(r, "advertiserId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламодателя")
		return
	}

	campaignID, err := uuid.Parse(chi.URLParam(r, "campaignId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID рекламной кампании")
		return
	}

	err = h.service.DeleteCampaign(ctx, advertiserID, campaignID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAdvertiserNotFound):
			WriteError(w, http.StatusNotFound, "Рекламодатель не найден", "")
		case errors.Is(err, domain.ErrAdNotFound):
			WriteError(w, http.StatusNotFound, "Рекламная кампания не найдена", "")
		default:
			log.Printf("[INTERNAL ERROR] failed to delete campaign: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GenerateAdText godoc
//
//	@Summary		Генерация текста рекламы
//	@Description	Генерирует текст рекламы на основе имени рекламодателя и названии рекламы
//	@Tags			Campaigns
//	@Accept			json
//	@Param			data	body	domain.GenerateAdTextRequest	true	"Информация для генерации текста"
//	@Produce		json
//	@Success		200	{object}	domain.GenerateAdTextResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/campaigns/generate [post]
func (h *CampaignHandler) GenerateAdText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var GenerateAdTextRequest *domain.GenerateAdTextRequest

	if err := json.NewDecoder(r.Body).Decode(&GenerateAdTextRequest); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	adText, err := h.service.GenerateAdText(ctx, GenerateAdTextRequest.AdvertiserName, GenerateAdTextRequest.AdTitle)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to generate ad text: %v", err)
		WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		return
	}

	response := domain.GenerateAdTextResponse{
		AdText: adText,
	}

	json.NewEncoder(w).Encode(response)
}

// SwitchModeration godoc
//
//	@Summary		Переключение модерации
//	@Description	Переключает статус модерации. Может быть true или false
//	@Tags			Campaigns
//	@Produce		json
//	@Success		200	{object}	domain.SwitchModerationResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/advertisers/campaigns/moderation [patch]
func (h *CampaignHandler) SwitchModeration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	newIsModerated, err := h.service.SwitchModeration(ctx)
	if err != nil {
		log.Printf("[INTERNAL ERROR] failed to switch moderation: %v", err)
		WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		return
	}

	response := domain.SwitchModerationResponse{
		IsModerated: newIsModerated,
	}

	json.NewEncoder(w).Encode(response)
}
