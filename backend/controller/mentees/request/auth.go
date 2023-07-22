package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type AuthMenteeInput struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (req *AuthMenteeInput) ToDomain() *mentees.MenteeAuth {
	return &mentees.MenteeAuth{
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
