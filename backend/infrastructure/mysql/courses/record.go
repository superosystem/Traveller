package courses

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/categories"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/mentors"
	"gorm.io/gorm"
)

type Course struct {
	ID          string `gorm:"primaryKey;size:200" json:"id"`
	MentorId    string `gorm:"size:200" json:"mentor_id"`
	CategoryId  string `gorm:"size:200" json:"category_id"`
	Title       string `gorm:"size:255" json:"title"`
	Description string `json:"description"`
	Thumbnail   string `gorm:"size:255" json:"thumbnail"`
	Category    categories.Category
	Mentor      mentors.Mentor
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (rec *Course) ToDomain() *courses.Domain {
	return &courses.Domain{
		ID:          rec.ID,
		MentorId:    rec.MentorId,
		CategoryId:  rec.CategoryId,
		Title:       rec.Title,
		Description: rec.Description,
		Thumbnail:   rec.Thumbnail,
		Category:    *rec.Category.ToDomain(),
		Mentor:      *rec.Mentor.ToDomain(),
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
		DeletedAt:   rec.DeletedAt,
	}
}

func FromDomain(courseDomain *courses.Domain) *Course {
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
	Category     categories.Category
	Mentor       mentors.Mentor
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
	TotalReviews int            `json:"total_reviews"`
	Rating       float32        `json:"rating"`
}

func (rec *CourseWithRating) ToDomain() *courses.Domain {
	return &courses.Domain{
		ID:           rec.ID,
		MentorId:     rec.MentorId,
		CategoryId:   rec.CategoryId,
		Title:        rec.Title,
		Description:  rec.Description,
		Thumbnail:    rec.Thumbnail,
		TotalReviews: rec.TotalReviews,
		Rating:       rec.Rating,
		Category:     *rec.Category.ToDomain(),
		Mentor:       *rec.Mentor.ToDomain(),
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
		DeletedAt:    rec.DeletedAt,
	}
}
