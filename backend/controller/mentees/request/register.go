package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type MenteeRegisterInput struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Phone    string `json:"phone" form:"phone" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	OTP      string `json:"otp" form:"otp" validate:"required"`
}

func (req *MenteeRegisterInput) ToDomain() *mentees.MenteeRegister {
	return &mentees.MenteeRegister{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
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
