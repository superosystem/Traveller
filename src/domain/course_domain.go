package domain

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID           string
	MentorId     string
	CategoryId   string
	Title        string
	Description  string
	Thumbnail    string
	TotalReviews int
	Rating       float32
	File         *multipart.FileHeader
	Category     Category
	Mentor       Mentor
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

type CourseRepository interface {
	// Create repository create course
	Create(courseDomain *Course) error

	// FindAll repository find all courses by course title and category
	FindAll(keyword string) (*[]Course, error)

	// FindById repository find course by id
	FindById(courseId string) (*Course, error)

	// FindByCategory repository find courses by category id
	FindByCategory(categoryId string) (*[]Course, error)

	// FindByMentor repository find courses by mentor id
	FindByMentor(mentorId string) (*[]Course, error)

	// FindByPopular repository find courses by highest rating
	FindByPopular() ([]Course, error)

	// Update repository update course
	Update(courseId string, courseDomain *Course) error

	// Delete repository delete course
	Delete(courseId string) error
}

type CourseUsecase interface {
	// Create usecase create new course
	Create(courseDomain *Course) error

	// FindAll usecase find all courses by course title and category
	FindAll(keyword string) (*[]Course, error)

	// FindById usecase find by id
	FindById(courseId string) (*Course, error)

	// FindByCategory usecase find by category id
	FindByCategory(categoryId string) (*[]Course, error)

	// FindByMentor usecase find courses by mentor id
	FindByMentor(mentorId string) (*[]Course, error)

	// FindByPopular usecase find courses by highest rating
	FindByPopular() ([]Course, error)

	// Update usecase update
	Update(courseId string, courseDomain *Course) error

	// Delete usecase delete
	Delete(courseId string) error
}
