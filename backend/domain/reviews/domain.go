package reviews

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type Domain struct {
	ID          string
	MenteeId    string
	CourseId    string
	Description string
	Rating      int
	Reviewed    bool
	Mentee      mentees.Domain
	Course      courses.Domain
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repository interface {
	// Create repository create new review course
	Create(reviewDomain *Domain) error

	// FindByCourse repository find all course reviews
	FindByCourse(courseId string) ([]Domain, error)
}

type Usecase interface {
	// Create usecase create new review course
	Create(reviewDomain *Domain) error

	// FindByCourse usecase find all course reviews
	FindByCourse(courseId string) ([]Domain, error)

	// FindByMentee usecase find all mentee reviews
	FindByMentee(menteeId string, title string) ([]Domain, error)
}
