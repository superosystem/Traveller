package entities

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MenteeProgress struct {
	ID         string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId   string `json:"mentee_id" gorm:"size:200"`
	CourseId   string `json:"course_id" gorm:"size:200"`
	MaterialId string `json:"material_id" gorm:"size:200"`
	Completed  string `json:"completed" gorm:"size:1"`
	Mentee     Mentee
	Course     Course
	Material   Material
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (req *MenteeProgress) ToMenteeProgressDomain() *domain.MenteeProgress {
	var completed bool

	if req.Completed == "1" {
		completed = true
	}

	return &domain.MenteeProgress{
		ID:         req.ID,
		MenteeId:   req.MenteeId,
		CourseId:   req.CourseId,
		MaterialId: req.MaterialId,
		Completed:  completed,
		Mentee:     *req.Mentee.ToMenteeDomain(),
		Course:     *req.Course.ToCourseDomain(),
		Material:   *req.Material.ToMaterialDomain(),
		CreatedAt:  req.CreatedAt,
		UpdatedAt:  req.UpdatedAt,
	}
}

func FromMenteeProgressDomain(menteeProgressDomain *domain.MenteeProgress) *MenteeProgress {
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
