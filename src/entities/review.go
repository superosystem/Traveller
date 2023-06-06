package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"time"
)

type Review struct {
	ID          string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId    string `json:"mentee_id" gorm:"size:200"`
	CourseId    string `json:"course_id" gorm:"size:200"`
	Description string `json:"description" gorm:"size:255"`
	Rating      int    `json:"rating"`
	Mentee      Mentee
	Course      Course
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (rec *Review) ToReviewDomain() *domain.Review {
	return &domain.Review{
		ID:          rec.ID,
		MenteeId:    rec.MenteeId,
		CourseId:    rec.CourseId,
		Description: rec.Description,
		Rating:      rec.Rating,
		Mentee:      *rec.Mentee.ToMenteeDomain(),
		Course:      *rec.Course.ToCourseDomain(),
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}

func FromReviewDomain(reviewDomain *domain.Review) *Review {
	return &Review{
		ID:          reviewDomain.ID,
		MenteeId:    reviewDomain.MenteeId,
		CourseId:    reviewDomain.CourseId,
		Description: reviewDomain.Description,
		Rating:      reviewDomain.Rating,
		CreatedAt:   reviewDomain.CreatedAt,
		UpdatedAt:   reviewDomain.UpdatedAt,
	}
}
