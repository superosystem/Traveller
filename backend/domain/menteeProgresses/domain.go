package mentee_progresses

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type Domain struct {
	ID            string
	MenteeId      string
	CourseId      string
	MaterialId    string
	ProgressCount int64
	Completed     bool
	Mentee        mentees.Domain
	Course        courses.Domain
	Material      materials.Domain
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Repository interface {
	// Add repository add new progress
	Add(menteeProgressDomain *Domain) error

	// FindByMaterial repository find progress by material
	FindByMaterial(menteeId string, materialId string) (*Domain, error)

	// FindByMentee repository find all progresses by mentee
	FindByMentee(menteeId string, courseId string) ([]Domain, error)

	// Count repository get mentee progresses count
	Count(menteeId string, title string, status string) ([]int64, error)

	// DeleteMenteeProgressesByCourse repository delete progress mentee by course
	DeleteMenteeProgressesByCourse(menteeId string, courseId string) error
}

type Usecase interface {
	// Add usecase add new progress
	Add(menteeProgressDomain *Domain) error

	// FindMaterialEnrolled usecase find material from enrolled course
	FindMaterialEnrolled(menteeId string, materialId string) (*Domain, error)
}
