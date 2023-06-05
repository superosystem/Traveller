package entities

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MenteeCourse struct {
	ID        string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId  string `json:"mentee_id" gorm:"size:200"`
	CourseId  string `json:"course_id" gorm:"size:200"`
	Status    string `json:"status" gorm:"size:50"`
	Reviewed  string `json:"reviewed" gorm:"size:1"`
	Mentee    Mentee
	Course    Course
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (rec *MenteeCourse) ToMenteeCourseDomain() *domain.MenteeCourse {
	reviewed := false

	if rec.Reviewed == "1" {
		reviewed = true
	}

	return &domain.MenteeCourse{
		ID:        rec.ID,
		MenteeId:  rec.MenteeId,
		CourseId:  rec.CourseId,
		Status:    rec.Status,
		Reviewed:  reviewed,
		Mentee:    *rec.Mentee.ToMenteeDomain(),
		Course:    *rec.Course.ToCourseDomain(),
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}

func FromMenteeCourseDomain(menteeCourseDomain *domain.MenteeCourse) *MenteeCourse {
	reviewed := "0"

	if menteeCourseDomain.Reviewed {
		reviewed = "1"
	}

	return &MenteeCourse{
		ID:        menteeCourseDomain.ID,
		MenteeId:  menteeCourseDomain.MenteeId,
		CourseId:  menteeCourseDomain.CourseId,
		Status:    menteeCourseDomain.Status,
		Reviewed:  reviewed,
		CreatedAt: menteeCourseDomain.CreatedAt,
		UpdatedAt: menteeCourseDomain.UpdatedAt,
	}
}
