package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var users []*domain.User

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newUsers, err := h.service.CreateUsers(ctx, users)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
		case domain.ErrBadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case domain.ErrInternalServerError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			log.Printf("[INTERNAL ERROR] failed to create client: %v", err)
			http.Error(w, domain.ErrInternalServerError.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUsers)
}