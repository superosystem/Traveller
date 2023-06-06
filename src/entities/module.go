package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"gorm.io/gorm"
	"time"
)

type Module struct {
	ID          string `json:"id" gorm:"primaryKey;size:200"`
	CourseId    string `json:"course_id" gorm:"size:200"`
	Title       string `json:"title" gorm:"size:255"`
	Description string `json:"description"`
	Course      Course
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (rec *Module) ToModuleDomain() *domain.Module {
	return &domain.Module{
		ID:          rec.ID,
		CourseId:    rec.CourseId,
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}

func FromModuleDomain(moduleDomain *domain.Module) *Module {
	return &Module{
		ID:          moduleDomain.ID,
		CourseId:    moduleDomain.CourseId,
		Title:       moduleDomain.Title,
		Description: moduleDomain.Description,
		CreatedAt:   moduleDomain.CreatedAt,
		UpdatedAt:   moduleDomain.UpdatedAt,
	}
}
