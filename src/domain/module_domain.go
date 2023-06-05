package domain

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	ID          string
	CourseId    string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type ModuleRepository interface {
	// Create repository create new module
	Create(moduleDomain *Module) error

	// FindById repository find module by id
	FindById(moduleId string) (*Module, error)

	// FindByCourse repository find modules by course
	FindByCourse(courseId string) ([]Module, error)

	// Update repository update module
	Update(moduleId string, moduleDomain *Module) error

	// Delete repository delete module
	Delete(moduleId string) error
}

type ModuleUsecase interface {
	// Create usecase create new module
	Create(moduleDomain *Module) error

	// FindById usecase find module by id
	FindById(moduleId string) (*Module, error)

	// Update usecase update module
	Update(moduleId string, moduleDomain *Module) error

	// Delete usecase delete module
	Delete(moduleId string) error
}
