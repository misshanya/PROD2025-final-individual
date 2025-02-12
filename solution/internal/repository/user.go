package repository

import (
	"context"

	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/infrastructure/db/sqlc/storage"
)

type UserRepository struct {
	queries *storage.Queries
}

func NewUserRepository(queries *storage.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) Create(ctx context.Context, users []*domain.User) ([]*domain.User, error) {
	for _, user := range users {
		err := r.queries.CreateUser(ctx, storage.CreateUserParams{
			ID: user.ID,
			Login: user.Login,
			Age: user.Age,
			Location: user.Location,
			Gender: user.Gender,
		})
		if err != nil {
			return []*domain.User{}, err
		}
	}
	
	return users, nil
}