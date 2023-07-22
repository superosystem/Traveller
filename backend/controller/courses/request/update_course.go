package request

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
)

type UpdateCourseInput struct {
	CategoryId  string                `json:"category_id" form:"category_id" validate:"required"`
	Title       string                `json:"title" form:"title" validate:"required"`
	Description string                `json:"description" form:"description" validate:"required"`
	Thumbnail   *multipart.FileHeader `json:"thumbnail,omitempty" form:"thumbnail"`
}

func (req *UpdateCourseInput) ToDomain() *courses.Domain {
	return &courses.Domain{
		CategoryId:  req.CategoryId,
		Title:       req.Title,
		Description: req.Description,
		File:        req.Thumbnail,
	}
}

func (req *UpdateCourseInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
