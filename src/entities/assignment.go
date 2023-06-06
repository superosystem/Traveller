package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	ID          string `gorm:"primaryKey;size:200" json:"id"`
	CourseId    string `json:"course_id" gorm:"size:200"`
	Course      Course
	Title       string         `gorm:"size:225" json:"title"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (rec *Assignment) ToAssignmentDomain() *domain.Assignment {
	return &domain.Assignment{
		ID:          rec.ID,
		CourseId:    rec.CourseId,
		Title:       rec.Title,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
		DeletedAt:   rec.DeletedAt,
	}
}

func FromAssignmentDomain(assignmentDomain *domain.Assignment) *Assignment {
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
