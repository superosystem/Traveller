package entities

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
	"gorm.io/gorm"
)

type Course struct {
	ID          string `gorm:"primaryKey;size:200" json:"id"`
	MentorId    string `gorm:"size:200" json:"mentor_id"`
	CategoryId  string `gorm:"size:200" json:"category_id"`
	Title       string `gorm:"size:255" json:"title"`
	Description string `json:"description"`
	Thumbnail   string `gorm:"size:255" json:"thumbnail"`
	Category    Category
	Mentor      Mentor
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (rec *Course) ToCourseDomain() *domain.Course {
	return &domain.Course{
		ID:          rec.ID,
		MentorId:    rec.MentorId,
		CategoryId:  rec.CategoryId,
		Title:       rec.Title,
		Description: rec.Description,
		Thumbnail:   rec.Thumbnail,
		Category:    *rec.Category.ToCategoryDomain(),
		Mentor:      *rec.Mentor.ToMentorDomain(),
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
		DeletedAt:   rec.DeletedAt,
	}
}

func FromCourseDomain(courseDomain *domain.Course) *Course {
	return &Course{
		ID:          courseDomain.ID,
		MentorId:    courseDomain.MentorId,
		CategoryId:  courseDomain.CategoryId,
		Title:       courseDomain.Title,
		Description: courseDomain.Description,
		Thumbnail:   courseDomain.Thumbnail,
		CreatedAt:   courseDomain.CreatedAt,
		UpdatedAt:   courseDomain.UpdatedAt,
		DeletedAt:   courseDomain.DeletedAt,
	}
}

type CourseWithRating struct {
	ID           string `json:"id"`
	MentorId     string `json:"mentor_id"`
	CategoryId   string `json:"category_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	Category     Category
	Mentor       Mentor
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
	TotalReviews int            `json:"total_reviews"`
	Rating       float32        `json:"rating"`
}

func (rec *CourseWithRating) ToDomain() *domain.Course {
	return &domain.Course{
		ID:           rec.ID,
		MentorId:     rec.MentorId,
		CategoryId:   rec.CategoryId,
		Title:        rec.Title,
		Description:  rec.Description,
		Thumbnail:    rec.Thumbnail,
		TotalReviews: rec.TotalReviews,
		Rating:       rec.Rating,
		Category:     *rec.Category.ToCategoryDomain(),
		Mentor:       *rec.Mentor.ToMentorDomain(),
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
		DeletedAt:    rec.DeletedAt,
	}
}
