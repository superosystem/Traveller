package domain

import "time"

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
	Add(menteeProgressDomain *MenteeProgress) error
	FindByMaterial(menteeId string, materialId string) (*MenteeProgress, error)
	FindByMentee(menteeId string, courseId string) ([]MenteeProgress, error)
	Count(menteeId string, title string, status string) ([]int64, error)
	DeleteMenteeProgressesByCourse(menteeId string, courseId string) error
}

type MenteeProgressUseCase interface {
	Add(menteeProgressDomain *MenteeProgress) error
	FindMaterialEnrolled(menteeId string, materialId string) (*MenteeProgress, error)
}
