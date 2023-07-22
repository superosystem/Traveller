package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
)

type AuthMentorInput struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (req *AuthMentorInput) ToDomain() *mentors.MentorAuth {
	return &mentors.MentorAuth{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *AuthMentorInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
