package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/TrainingSystem/backend/domain/categories"
)

type Category struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func (req *Category) ToDomain() *categories.Domain {
	return &categories.Domain{
		Name: req.Name,
	}
}

func (req *Category) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
