package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type OTP struct {
	Key string `json:"email" form:"email" validate:"required,email"`
}

func (req *OTP) ToDomain() *domain.Otp {
	return &domain.Otp{
		Key: req.Key,
	}
}

func (req *OTP) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
