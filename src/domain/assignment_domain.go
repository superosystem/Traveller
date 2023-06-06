package domain

import (
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	ID          string
	CourseId    string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type AssignmentRepository interface {
	Create(assignmentDomain *Assignment) error
	FindById(assignmentId string) (*Assignment, error)
	FindByCourseId(courseId string) (*Assignment, error)
	FindByCourses(courseIds []string) (*[]Assignment, error)
	Update(assignmentId string, assignmentDomain *Assignment) error
	Delete(assignmentId string) error
}

type AssignmentUseCase interface {
	Create(assignmentDomain *Assignment) error
	FindById(assignmentId string) (*Assignment, error)
	FindByCourseId(courseId string) (*Assignment, error)
	Update(assignmentId string, assignmentDomain *Assignment) error
	Delete(assignmentId string) error
}
