package domain

import (
	"gorm.io/gorm"
	"time"
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
	Create(moduleDomain *Module) error
	FindById(moduleId string) (*Module, error)
	FindByCourse(courseId string) ([]Module, error)
	Update(moduleId string, moduleDomain *Module) error
	Delete(moduleId string) error
}

type ModuleUseCase interface {
	Create(moduleDomain *Module) error
	FindById(moduleId string) (*Module, error)
	Update(moduleId string, moduleDomain *Module) error
	Delete(moduleId string) error
}
