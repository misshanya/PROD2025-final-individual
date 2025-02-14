package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gitlab.prodcontest.ru/2025-final-projects-back/misshanya/internal/domain"
)

type TimeRepository struct {
	rdb *redis.Client
}

func NewTimeRepository(rdb *redis.Client) *TimeRepository {
	return &TimeRepository{
		rdb: rdb,
	}
}

func (r *TimeRepository) SetCurrentDate(ctx context.Context, newDate int) error {
	currentDateStr, err := r.rdb.Get(ctx, "current_date").Result()
	if err != nil && err != redis.Nil {
		return err
	} else if err == redis.Nil {
		currentDateStr = "0"
	}
	currentDate, err := strconv.Atoi(currentDateStr)
	if err != nil {
		return err
	}
	if newDate < currentDate {
		return domain.ErrNewDateLowerThanCurrent
	}

	err = r.rdb.Set(ctx, "current_date", newDate, 0).Err()
	return err
}