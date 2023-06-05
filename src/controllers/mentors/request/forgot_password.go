package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type ForgotPasswordInput struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

func (req *ForgotPasswordInput) ToDomain() *domain.MentorForgotPassword {
	return &domain.MentorForgotPassword{
		Email: req.Email,
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
