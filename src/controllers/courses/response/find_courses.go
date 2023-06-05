package response

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type FindCourses struct {
	CourseId     string    `json:"course_id"`
	MentorId     string    `json:"mentor_id"`
	CategoryId   string    `json:"category_id"`
	Mentor       string    `json:"mentor"`
	Category     string    `json:"category"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Thumbnail    string    `json:"thumbnail"`
	TotalReviews int       `json:"total_reviews"`
	Rating       float32   `json:"rating"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func AllCourses(res *domain.Course) FindCourses {
	return FindCourses{
		CourseId:     res.ID,
		MentorId:     res.MentorId,
		CategoryId:   res.CategoryId,
		Mentor:       res.Mentor.Fullname,
		Category:     res.Category.Name,
		Title:        res.Title,
		Description:  res.Description,
		Thumbnail:    res.Thumbnail,
		TotalReviews: res.TotalReviews,
		Rating:       res.Rating,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
}
