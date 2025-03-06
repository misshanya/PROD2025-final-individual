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

type UserHandler struct {
	service *app.UserService
}

func NewUserHandler(service *app.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// CreateUsers godoc
//
//	@Summary		Массовое создание/обновление клиентов
//	@Description	Создает новых или обновляет существующих клиентов
//	@Tags			Clients
//	@Accept			json
//	@Param			clients	body	[]domain.User	true	"Clients"
//	@Produce		json
//	@Success		201	{object}	[]domain.User
//	@Failure		400 {object} ErrorResponse
//	@Failure		500	{object} ErrorResponse
//	@Router			/clients/bulk [post]
func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var users []*domain.User

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", err.Error())
		return
	}

	newUsers, err := h.service.CreateUpdateUsers(ctx, users)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			WriteError(w, http.StatusBadRequest, "Некорректный запрос", "")
		default:
			log.Printf("[INTERNAL ERROR] failed to create client: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUsers)
}

// GetByID godoc
//
//	@Summary		Получение клиента по ID
//	@Description	Возвращает информацию о клиенте по его ID
//	@Tags			Clients
//	@Param			clientId	path	string	true	"UUID клиента"
//	@Produce		json
//	@Success		200	{object}	[]domain.User
//	@Failure		400	{object} ErrorResponse
//	@Failure		404	{object} ErrorResponse
//	@Failure		500	{object} ErrorResponse
//	@Router			/clients/{clientId} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clientID, err := uuid.Parse(chi.URLParam(r, "clientId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Некорректный запрос", "невалидный ID клиента")
		return
	}

	client, err := h.service.GetByID(ctx, clientID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			WriteError(w, http.StatusNotFound, "Пользователь не найден", "")
			return
		default:
			log.Printf("[INTERNAL ERROR] failed to get client: %v", err)
			WriteError(w, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), "")
		}
	}

	json.NewEncoder(w).Encode(client)
}
