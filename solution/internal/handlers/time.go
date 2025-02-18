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

// SetCurrentDate godoc
//
//	@Summary		Установка текущей даты
//	@Description	Устанавливает текущий день в системе в заданную дату
//	@Tags			Time
//	@Accept			json
//	@Param			newDate	body	domain.CurrentDate	true	"Новый текущий день"
//	@Produce		json
//	@Success		200	{object}	domain.CurrentDate
//	@Failure		400	{object} ErrorResponse
//	@Failure		500	{object} ErrorResponse
//	@Router			/time/advance [post]
func (h *TimeHandler) SetCurrentDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body domain.CurrentDate

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	if err := h.service.SetCurrentDate(ctx, body.CurrentDate); err != nil {
		switch err {
		case domain.ErrNewDateLowerThanCurrent:
			WriteError(w, http.StatusBadRequest, "Некорректный запрос", "новая дата меньше текущей")
			return
		default:
			log.Printf("[INTERNAL ERROR] failed to set current date: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
			return
		}
	}

	json.NewEncoder(w).Encode(body)
}
