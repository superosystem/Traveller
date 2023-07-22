package response

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/reviews"
)

type FindReviewByMentee struct {
	MenteeId  string    `json:"mentee_id"`
	CourseId  string    `json:"course_id"`
	Title     string    `json:"title"`
	Mentor    string    `json:"mentor"`
	Reviewed  bool      `json:"reviewed"`
	Thumbnail string    `json:"thumbnail"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ReviewsByMentee(reviewDomain *reviews.Domain) FindReviewByMentee {
	return FindReviewByMentee{
		MenteeId:  reviewDomain.MenteeId,
		CourseId:  reviewDomain.CourseId,
		Title:     reviewDomain.Course.Title,
		Mentor:    reviewDomain.Course.Mentor.Fullname,
		Reviewed:  reviewDomain.Reviewed,
		Thumbnail: reviewDomain.Course.Thumbnail,
		CreatedAt: reviewDomain.CreatedAt,
		UpdatedAt: reviewDomain.UpdatedAt,
	}
}
