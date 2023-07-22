package mentee_progresses

import (
	"time"

	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/courses"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/materials"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/mentees"
)

type MenteeProgress struct {
	ID         string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId   string `json:"mentee_id" gorm:"size:200"`
	CourseId   string `json:"course_id" gorm:"size:200"`
	MaterialId string `json:"material_id" gorm:"size:200"`
	Completed  string `json:"completed" gorm:"size:1"`
	Mentee     mentees.Mentee
	Course     courses.Course
	Material   materials.Material
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (req *MenteeProgress) ToDomain() *menteeProgresses.Domain {
	var completed bool

	if req.Completed == "1" {
		completed = true
	}

	return &menteeProgresses.Domain{
		ID:         req.ID,
		MenteeId:   req.MenteeId,
		CourseId:   req.CourseId,
		MaterialId: req.MaterialId,
		Completed:  completed,
		Mentee:     *req.Mentee.ToDomain(),
		Course:     *req.Course.ToDomain(),
		Material:   *req.Material.ToDomain(),
		CreatedAt:  req.CreatedAt,
		UpdatedAt:  req.UpdatedAt,
	}
}

func FromDomain(menteeProgressDomain *menteeProgresses.Domain) *MenteeProgress {
	var completed string

	if menteeProgressDomain.Completed {
		completed = "1"
	}

	return &MenteeProgress{
		ID:         menteeProgressDomain.ID,
		MenteeId:   menteeProgressDomain.MenteeId,
		CourseId:   menteeProgressDomain.CourseId,
		MaterialId: menteeProgressDomain.MaterialId,
		Completed:  completed,
		CreatedAt:  menteeProgressDomain.CreatedAt,
		UpdatedAt:  menteeProgressDomain.UpdatedAt,
	}
}
