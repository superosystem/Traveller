package request

import (
	"github.com/go-playground/validator/v10"
	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
)

type AddProgressInput struct {
	MenteeId   string `json:"mentee_id" form:"mentee_id"`
	CourseId   string `json:"course_id" form:"course_id"`
	MaterialId string `json:"material_id" form:"material_id"`
}

func (req *AddProgressInput) ToDomain() *menteeProgresses.Domain {
	return &menteeProgresses.Domain{
		MenteeId:   req.MenteeId,
		CourseId:   req.CourseId,
		MaterialId: req.MaterialId,
	}
}

func (req *AddProgressInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
