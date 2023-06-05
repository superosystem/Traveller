package domain

import (
	"time"

	"gorm.io/gorm"
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
	// Create repository create new assignment
	Create(assignmentDomain *Assignment) error

	// FindById repository find assignment by id
	FindById(assignmentId string) (*Assignment, error)

	// FindByCourseId repository find assignment by courseid
	FindByCourseId(courseId string) (*Assignment, error)

	// FindByCourses repository find assignments by courses
	FindByCourses(courseIds []string) (*[]Assignment, error)

	// Update repository update assignment
	Update(assignmentId string, assignmentDomain *Assignment) error

	// Delete repository delete assignment
	Delete(assignmentId string) error
}

type AssignmentUsecase interface {
	// Create usecase create new assignment
	Create(assignmentDomain *Assignment) error

	// FindById usecase findfind assignment by id
	FindById(assignmentId string) (*Assignment, error)

	// FindByCourseId usecase find assignment by courseid
	FindByCourseId(courseId string) (*Assignment, error)

	// Update usecase update assignment
	Update(assignmentId string, assignmentDomain *Assignment) error

	// Delete usecase delete assignment
	Delete(assignmentId string) error
}
