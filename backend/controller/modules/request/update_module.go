package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/modules"
)

type UpdateModuleInput struct {
	CourseId    string `json:"course_id" form:"course_id" validate:"required"`
	Title       string `json:"title" form:"title" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

func (req *UpdateModuleInput) ToDomain() *modules.Domain {
	return &modules.Domain{
		CourseId:    req.CourseId,
		Title:       req.Title,
		Description: req.Description,
	}
}

func (req *UpdateModuleInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
