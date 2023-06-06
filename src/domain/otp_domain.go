package domain

import (
	"context"
	"time"
)

type Otp struct {
	Key   string
	Value string
}

type OtpRepository interface {
	Save(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type OtpUseCase interface {
	SendOTP(otpDomain *Otp) error
	CheckOTP(otpDomain *Otp) error
}
