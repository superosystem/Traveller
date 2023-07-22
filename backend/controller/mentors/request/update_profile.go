package request

import (
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
)

type MentorUpdateProfile struct {
	UserID             string                `json:"user_id" form:"user_id" validate:"required"`
	Fullname           string                `json:"fullname" form:"fullname"`
	Email              string                `json:"email" form:"email"`
	Phone              string                `json:"phone" form:"phone"`
	Jobs               string                `json:"jobs" form:"jobs" `
	Gender             string                `json:"gender" form:"gender"`
	BirthPlace         string                `json:"birth_place" form:"birth_place"`
	BirthDate          string                `json:"birth_date" form:"birth_date"`
	Address            string                `json:"address" form:"address"`
	ProfilePictureFile *multipart.FileHeader `json:"profile_picture" form:"profile_picture"`
}

func (req *MentorUpdateProfile) ToDomain() *mentors.MentorUpdateProfile {

	birth, _ := time.Parse("2006-01-02", req.BirthDate)

	return &mentors.MentorUpdateProfile{

		UserID:             req.UserID,
		Fullname:           req.Fullname,
		Email:              req.Email,
		Phone:              req.Phone,
		Jobs:               req.Jobs,
		Gender:             req.Gender,
		BirthPlace:         req.BirthPlace,
		BirthDate:          birth,
		Address:            req.Address,
		ProfilePictureFile: req.ProfilePictureFile,
	}
}

func (req *MentorUpdateProfile) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
