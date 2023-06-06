package domain

import "time"

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
	Create(reviewDomain *Review) error
	FindByCourse(courseId string) ([]Review, error)
}

type ReviewUseCase interface {
	Create(reviewDomain *Review) error
	FindByCourse(courseId string) ([]Review, error)
	FindByMentee(menteeId string, title string) ([]Review, error)
}
