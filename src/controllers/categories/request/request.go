package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type Category struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func (req *Category) ToDomain() *domain.Category {
	return &domain.Category{
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
