package request

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type MenteeUpdateProfile struct {
	Fullname           string                `json:"fullname" form:"fullname"`
	Phone              string                `json:"phone" form:"phone"`
	BirthDate          string                `json:"birth_date" form:"birth_date"`
	ProfilePictureFile *multipart.FileHeader `json:"profile_picture" form:"profile_picture"`
}

func (req *MenteeUpdateProfile) ToDomain() *mentees.Domain {
	return &mentees.Domain{
		Fullname:           req.Fullname,
		Phone:              req.Phone,
		BirthDate:          req.BirthDate,
		ProfilePictureFile: req.ProfilePictureFile,
	}
}

func (req *MenteeUpdateProfile) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
