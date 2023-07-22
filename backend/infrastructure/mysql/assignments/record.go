package assignments

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/assignments"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/courses"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          string `gorm:"primaryKey;size:200" json:"id"`
	CourseId    string `json:"course_id" gorm:"size:200"`
	Course      courses.Course
	Title       string         `gorm:"size:225" json:"title"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (rec *Assignment) ToDomain() *assignments.Domain {
	return &assignments.Domain{
		ID:          rec.ID,
		CourseId:    rec.CourseId,
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
		DeletedAt:   rec.DeletedAt,
	}
}

func FromDomain(assignmentDomain *assignments.Domain) *Assignment {
	return &Assignment{
		ID:          assignmentDomain.ID,
		CourseId:    assignmentDomain.CourseId,
		Title:       assignmentDomain.Title,
		Description: assignmentDomain.Description,
		CreatedAt:   assignmentDomain.CreatedAt,
		UpdatedAt:   assignmentDomain.UpdatedAt,
		DeletedAt:   assignmentDomain.DeletedAt,
	}
}
