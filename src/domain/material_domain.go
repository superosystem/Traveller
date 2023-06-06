package domain

import (
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

type Material struct {
	ID          string
	CourseId    string
	ModuleId    string
	Title       string
	URL         string
	Description string
	File        *multipart.FileHeader
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type MaterialRepository interface {
	Create(materialDomain *Material) error
	FindById(materialId string) (*Material, error)
	FindByModule(moduleIds []string) ([]Material, error)
	CountByCourse(courseIds []string) ([]int64, error)
	Update(materialId string, materialDomain *Material) error
	Delete(materialId string) error
	Deletes(moduleId string) error
}

type MaterialUseCase interface {
	Create(materialDomain *Material) error
	FindById(materialId string) (*Material, error)
	Update(materialId string, materialDomain *Material) error
	Delete(materialId string) error
	Deletes(moduleId string) error
}
