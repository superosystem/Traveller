package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type CheckOTP struct {
	// the email
	Key string `json:"email" form:"email" validate:"required,email"`

	// the OTP
	Value string `json:"otp" form:"otp" validate:"required"`
}

func (req *CheckOTP) ToDomain() *domain.Otp {
	return &domain.Otp{
		Key:   req.Key,
		Value: req.Value,
	}
}

func (req *CheckOTP) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return err
}
