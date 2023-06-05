package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type EnrollCourse struct {
	MenteeId string `json:"mentee_id" form:"mentee_id" validate:"required"`
	CourseId string `json:"course_id" form:"course_id" validate:"required"`
}

func (req *EnrollCourse) ToDomain() *domain.MenteeCourse {
	return &domain.MenteeCourse{
		MenteeId: req.MenteeId,
		CourseId: req.CourseId,
	}
}

func (req *EnrollCourse) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
