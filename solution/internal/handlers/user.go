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
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/clients/bulk [post]
func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var users []*domain.User

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newUsers, err := h.service.CreateUpdateUsers(ctx, users)
	if err != nil {
		switch err {
		case domain.ErrBadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Printf("[INTERNAL ERROR] failed to create client: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
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
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		404	{string}	string	"Client not found"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/clients/{clientId} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clientID, err := uuid.Parse(chi.URLParam(r, "clientId"))
	if err != nil {
		log.Printf("failed to convert clientID to uuid: %v", err)
		http.Error(w, domain.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	client, err := h.service.GetByID(ctx, clientID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			log.Printf("[INTERNAL ERROR] failed to get client: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		}
	}

	json.NewEncoder(w).Encode(client)
}
