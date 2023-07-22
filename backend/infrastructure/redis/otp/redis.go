package otp

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/superosystem/TrainingSystem/backend/domain/otp"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type otpRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) otp.Repository {
	return otpRepository{
		client: client,
	}
}

func (o otpRepository) Save(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := o.client.Set(ctx, key, value, ttl).Err()

	if err != nil {
		return err
	}

	return nil
}

func (o otpRepository) Get(ctx context.Context, key string) (string, error) {
	result, err := o.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", helper.ErrOTPExpired
	} else if err != nil {
		return "", err
	}

	return result, nil
}
