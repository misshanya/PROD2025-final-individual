package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type MLRepository struct {
	rdb *redis.Client
}

func NewMLRepository(rdb *redis.Client) *MLRepository {
	return &MLRepository{
		rdb: rdb,
	}
}

func (r *MLRepository) SwitchModeration(ctx context.Context) (bool, error) {
	var isModerated bool
	isModeratedStr, err := r.rdb.Get(ctx, "isModerated").Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	isModerated = isModeratedStr == "1"

	err = r.rdb.Set(ctx, "isModerated", !isModerated, 0).Err()
	if err != nil {
		return false, err
	}

	return !isModerated, nil
}

func (r *MLRepository) CheckModeration(ctx context.Context) (bool, error) {
	var isModerated bool
	isModeratedStr, err := r.rdb.Get(ctx, "isModerated").Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	isModerated = isModeratedStr == "1"
	return isModerated, nil
}
