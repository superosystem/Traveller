package response

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/reviews"
)

type FindReviewByCourse struct {
	ID             string    `json:"id"`
	MenteeId       string    `json:"mentee_id"`
	CourseId       string    `json:"course_id"`
	Mentee         string    `json:"mentee"`
	ProfilePicture string    `json:"profile_picture"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Rating         int       `json:"rating"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ReviewsByCourse(reviewDomain *reviews.Domain) FindReviewByCourse {
	return FindReviewByCourse{
		ID:             reviewDomain.ID,
		MenteeId:       reviewDomain.MenteeId,
		CourseId:       reviewDomain.CourseId,
		Mentee:         reviewDomain.Mentee.Fullname,
		ProfilePicture: reviewDomain.Mentee.ProfilePicture,
		Title:          reviewDomain.Course.Title,
		Description:    reviewDomain.Description,
		Rating:         reviewDomain.Rating,
		CreatedAt:      reviewDomain.CreatedAt,
		UpdatedAt:      reviewDomain.UpdatedAt,
	}
}
