package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"time"
)

type otpRepository struct {
	client *redis.Client
}

func NewOtpRepository(client *redis.Client) domain.OtpRepository {
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
