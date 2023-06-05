package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type ForgotPasswordInput struct {
	Email            string `json:"email" form:"email" validate:"required,email"`
	Password         string `json:"password" form:"password" validate:"required"`
	RepeatedPassword string `json:"repeated_password" form:"repeated_password" validate:"required"`
	OTP              string `json:"otp" form:"otp" validate:"required"`
}

func (req *ForgotPasswordInput) ToDomain() *domain.MenteeForgotPassword {
	return &domain.MenteeForgotPassword{
		Email:            req.Email,
		Password:         req.Password,
		RepeatedPassword: req.RepeatedPassword,
		OTP:              req.OTP,
	}
}

func (req *ForgotPasswordInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
