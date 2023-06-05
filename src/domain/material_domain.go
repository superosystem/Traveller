package domain

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
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
	// Create repository create materials
	Create(materialDomain *Material) error

	// FindById repository find material by id
	FindById(materialId string) (*Material, error)

	// FindByModule repository find materials by module
	FindByModule(moduleIds []string) ([]Material, error)

	// CountByCourse repository find total materials by course
	CountByCourse(courseIds []string) ([]int64, error)

	// Update repository update material
	Update(materialId string, materialDomain *Material) error

	// Delete repository delete single material by id material
	Delete(materialId string) error

	// Deletes repository delete materials by id module
	Deletes(moduleId string) error
}

type MaterialUsecase interface {
	// Create usecase create material
	Create(materialDomain *Material) error

	// FindById usecase find material by id
	FindById(materialId string) (*Material, error)

	// Update usecase update material
	Update(materialId string, materialDomain *Material) error

	// Delete usecase detele material by id material
	Delete(materialId string) error

	// Deletes usecase delete materials by id module
	Deletes(moduleId string) error
}
