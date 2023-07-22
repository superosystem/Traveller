package modules

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/modules"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/courses"
	"gorm.io/gorm"
)

type Module struct {
	ID          string `json:"id" gorm:"primaryKey;size:200"`
	CourseId    string `json:"course_id" gorm:"size:200"`
	Title       string `json:"title" gorm:"size:255"`
	Description string `json:"description"`
	Course      courses.Course
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (rec *Module) ToDomain() *modules.Domain {
	return &modules.Domain{
		ID:          rec.ID,
		CourseId:    rec.CourseId,
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}

func FromDomain(moduleDomain *modules.Domain) *Module {
	return &Module{
		ID:          moduleDomain.ID,
		CourseId:    moduleDomain.CourseId,
		Title:       moduleDomain.Title,
		Description: moduleDomain.Description,
		CreatedAt:   moduleDomain.CreatedAt,
		UpdatedAt:   moduleDomain.UpdatedAt,
	}
}
