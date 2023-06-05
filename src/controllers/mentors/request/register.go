package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MentorRegisterInput struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (req *MentorRegisterInput) ToDomain() *domain.MentorRegister {
	return &domain.MentorRegister{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *MentorRegisterInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
