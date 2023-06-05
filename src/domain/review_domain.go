package domain

import (
	"time"
)

type Review struct {
	ID          string
	MenteeId    string
	CourseId    string
	Description string
	Rating      int
	Reviewed    bool
	Mentee      Mentee
	Course      Course
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReviewRepository interface {
	// Create repository create new review course
	Create(reviewDomain *Review) error

	// FindByCourse repository find all course reviews
	FindByCourse(courseId string) ([]Review, error)
}

type ReviewUsecase interface {
	// Create usecase create new review course
	Create(reviewDomain *Review) error

	// FindByCourse usecase find all course reviews
	FindByCourse(courseId string) ([]Review, error)

	// FindByMentee usecase find all mentee reviews
	FindByMentee(menteeId string, title string) ([]Review, error)
}
