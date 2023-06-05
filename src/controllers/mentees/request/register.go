package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MenteeRegisterInput struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Phone    string `json:"phone" form:"phone" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
	OTP      string `json:"otp" form:"otp" validate:"required"`
}

func (req *MenteeRegisterInput) ToDomain() *domain.MenteeRegister {
	return &domain.MenteeRegister{
		Fullname: req.Fullname,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
		OTP:      req.OTP,
	}
}

func (req *MenteeRegisterInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
