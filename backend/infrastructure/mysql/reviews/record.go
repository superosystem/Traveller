package reviews

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/reviews"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/courses"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/mentees"
)

type Review struct {
	ID          string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId    string `json:"mentee_id" gorm:"size:200"`
	CourseId    string `json:"course_id" gorm:"size:200"`
	Description string `json:"description" gorm:"size:255"`
	Rating      int    `json:"rating"`
	Mentee      mentees.Mentee
	Course      courses.Course
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (rec *Review) ToDomain() *reviews.Domain {
	return &reviews.Domain{
		ID:          rec.ID,
		MenteeId:    rec.MenteeId,
		CourseId:    rec.CourseId,
		Description: rec.Description,
		Rating:      rec.Rating,
		Mentee:      *rec.Mentee.ToDomain(),
		Course:      *rec.Course.ToDomain(),
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}

func FromDomain(reviewDomain *reviews.Domain) *Review {
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
