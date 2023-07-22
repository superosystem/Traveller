package otp

import (
	"context"
	"time"
)

type Domain struct {
	Key   string
	Value string
}

type Repository interface {
	Save(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Usecase interface {
	SendOTP(otpDomain *Domain) error
	CheckOTP(otpDomain *Domain) error
}
