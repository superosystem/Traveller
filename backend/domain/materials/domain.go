package materials

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Domain struct {
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

type Repository interface {
	// Create repository create materials
	Create(materialDomain *Domain) error

	// FindById repository find material by id
	FindById(materialId string) (*Domain, error)

	// FindByModule repository find materials by module
	FindByModule(moduleIds []string) ([]Domain, error)

	// CountByCourse repository find total materials by course
	CountByCourse(courseIds []string) ([]int64, error)

	// Update repository update material
	Update(materialId string, materialDomain *Domain) error

	// Delete repository delete single material by id material
	Delete(materialId string) error

	// Deletes repository delete materials by id module
	Deletes(moduleId string) error
}

type Usecase interface {
	// Create usecase create material
	Create(materialDomain *Domain) error

	// FindById usecase find material by id
	FindById(materialId string) (*Domain, error)

	// Update usecase update material
	Update(materialId string, materialDomain *Domain) error

	// Delete usecase detele material by id material
	Delete(materialId string) error

	// Deletes usecase delete materials by id module
	Deletes(moduleId string) error
}
