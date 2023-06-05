package request

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type UpdateMaterialInput struct {
	ModuleId    string                `json:"module_id" form:"module_id" validate:"required"`
	Title       string                `json:"title" form:"title" validate:"required"`
	Description string                `json:"description" form:"description" validate:"required"`
	File        *multipart.FileHeader `json:"video" form:"video"`
}

func (req *UpdateMaterialInput) ToDomain() *domain.Material {
	return &domain.Material{
		ModuleId:    req.ModuleId,
		Title:       req.Title,
		Description: req.Description,
		File:        req.File,
	}
}

func (req *UpdateMaterialInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
