package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		_, err := r.queries.GetUserByID(ctx, user.ID)
		if err == nil {
			return []*domain.User{}, domain.ErrUserAlreadyExists
		}

		_, err = r.queries.CreateUser(ctx, storage.CreateUserParams{
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

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &domain.User{
		ID: user.ID,
		Login: user.Login,
		Age: user.Age,
		Location: user.Location,
		Gender: user.Gender,
	}, nil
}