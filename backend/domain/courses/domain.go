package courses

import (
	"mime/multipart"
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/categories"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
	"gorm.io/gorm"
)

type Domain struct {
	ID           string
	MentorId     string
	CategoryId   string
	Title        string
	Description  string
	Thumbnail    string
	TotalReviews int
	Rating       float32
	File         *multipart.FileHeader
	Category     categories.Domain
	Mentor       mentors.Domain
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

type Repository interface {
	// Create repository create course
	Create(courseDomain *Domain) error

	// FindAll repository find all courses by course title and category
	FindAll(keyword string) (*[]Domain, error)

	// FindById repository find course by id
	FindById(courseId string) (*Domain, error)

	// FindByCategory repository find courses by category id
	FindByCategory(categoryId string) (*[]Domain, error)

	// FindByMentor repository find courses by mentor id
	FindByMentor(mentorId string) (*[]Domain, error)

	// FindByPopular repository find courses by highest rating
	FindByPopular() ([]Domain, error)

	// Update repository update course
	Update(courseId string, courseDomain *Domain) error

	// Delete repository delete course
	Delete(courseId string) error
}

type Usecase interface {
	// Create usecase create new course
	Create(courseDomain *Domain) error

	// FindAll usecase find all courses by course title and category
	FindAll(keyword string) (*[]Domain, error)

	// FindById usecase find by id
	FindById(courseId string) (*Domain, error)

	// FindByCategory usecase find by category id
	FindByCategory(categoryId string) (*[]Domain, error)

	// FindByMentor usecase find courses by mentor id
	FindByMentor(mentorId string) (*[]Domain, error)

	// FindByPopular usecase find courses by highest rating
	FindByPopular() ([]Domain, error)

	// Update usecase update
	Update(courseId string, courseDomain *Domain) error

	// Delete usecase delete
	Delete(courseId string) error
}
