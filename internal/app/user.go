package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUpdateUsers(ctx context.Context, users []*domain.User) ([]*domain.User, error) {
	for _, user := range users {
		if err := validateUser(user); err != nil {
			return []*domain.User{}, domain.ErrBadRequest
		}
	}

	newUsers, err := s.repo.CreateUpdateUsers(ctx, users)
	if err != nil {
		return []*domain.User{}, err
	}

	return newUsers, nil
}

func validateUser(user *domain.User) error {
	if user.Age <= 0 || user.Age >= 200 {
		return fmt.Errorf("age must be greater than 0 and lower than 200")
	}
	if !isValidGender(user.Gender) {
		return fmt.Errorf("invalid user gender")
	}
	return nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
