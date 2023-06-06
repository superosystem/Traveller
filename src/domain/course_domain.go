package domain

import (
	"gorm.io/gorm"
	"mime/multipart"
	"time"
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
	Create(courseDomain *Course) error
	FindAll(keyword string) (*[]Course, error)
	FindById(courseId string) (*Course, error)
	FindByCategory(categoryId string) (*[]Course, error)
	FindByMentor(mentorId string) (*[]Course, error)
	FindByPopular() ([]Course, error)
	Update(courseId string, courseDomain *Course) error
	Delete(courseId string) error
}

type CourseUseCase interface {
	Create(courseDomain *Course) error
	FindAll(keyword string) (*[]Course, error)
	FindById(courseId string) (*Course, error)
	FindByCategory(categoryId string) (*[]Course, error)
	FindByMentor(mentorId string) (*[]Course, error)
	FindByPopular() ([]Course, error)
	Update(courseId string, courseDomain *Course) error
	Delete(courseId string) error
}
