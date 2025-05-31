package tempid_repo

import (
	"context"
	"errors"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/redis/go-redis/v9"
)

type TempIdRepo struct {
	redisClient *redis.Client
}

func NewTempIdRepo(redisClient *redis.Client) *TempIdRepo {
	return &TempIdRepo{
		redisClient: redisClient,
	}
}

func (t *TempIdRepo) Get(ctx context.Context, tempId string) (string, error) {
	result, err := t.redisClient.Get(ctx, tempId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", apperrors.NotFound()
		}
		return "", err
	}
	return result, nil
}

func (t *TempIdRepo) Set(ctx context.Context, tempId string, data string, exp time.Duration) error {
	return t.redisClient.Set(ctx, tempId, data, exp).Err()
}

func (t *TempIdRepo) Delete(ctx context.Context, tempId string) error {
	err := t.redisClient.Del(ctx, tempId).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return apperrors.NotFound()
		}
		return err
	}

	return nil
}
