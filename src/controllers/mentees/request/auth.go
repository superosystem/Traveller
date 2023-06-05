package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type AuthMenteeInput struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (req *AuthMenteeInput) ToDomain() *domain.MenteeAuth {
	return &domain.MenteeAuth{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *AuthMenteeInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
