package domain

import (
	"time"
)

type MenteeProgress struct {
	ID            string
	MenteeId      string
	CourseId      string
	MaterialId    string
	ProgressCount int64
	Completed     bool
	Mentee        Mentee
	Course        Course
	Material      Material
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type MenteeProgressRepository interface {
	// Add repository add new progress
	Add(menteeProgressDomain *MenteeProgress) error

	// FindByMaterial repository find progress by material
	FindByMaterial(menteeId string, materialId string) (*MenteeProgress, error)

	// FindByMentee repository find all progresses by mentee
	FindByMentee(menteeId string, courseId string) ([]MenteeProgress, error)

	// Count repository get mentee progresses count
	Count(menteeId string, title string, status string) ([]int64, error)

	// DeleteMenteeProgressesByCourse repository delete progress mentee by course
	DeleteMenteeProgressesByCourse(menteeId string, courseId string) error
}

type MenteeProgressUsecase interface {
	// Add usecase add new progress
	Add(menteeProgressDomain *MenteeProgress) error

	// FindMaterialEnrolled usecase find material from enrolled course
	FindMaterialEnrolled(menteeId string, materialId string) (*MenteeProgress, error)
}
