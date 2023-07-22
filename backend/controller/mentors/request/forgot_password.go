package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
)

type ForgotPasswordInput struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

func (req *ForgotPasswordInput) ToDomain() *mentors.MentorForgotPassword {
	return &mentors.MentorForgotPassword{
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
