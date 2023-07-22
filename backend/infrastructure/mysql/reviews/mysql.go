package reviews

import (
	"github.com/superosystem/TrainingSystem/backend/domain/reviews"
	"gorm.io/gorm"
)

type reviewRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) reviews.Repository {
	return reviewRepository{
		conn: conn,
	}
}

func (rr reviewRepository) Create(reviewDomain *reviews.Domain) error {
	rec := FromDomain(reviewDomain)

	err := rr.conn.Model(&Review{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (rr reviewRepository) FindByCourse(courseId string) ([]reviews.Domain, error) {
	var rec []Review

	err := rr.conn.Model(&Review{}).Preload("Mentee").Preload("Course").
		Joins("INNER JOIN mentees ON mentees.id = reviews.mentee_id").
		Joins("INNER JOIN courses ON courses.id = reviews.course_id").
		Where("reviews.course_id = ?", courseId).
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var reviewDomain []reviews.Domain

	for _, review := range rec {
		reviewDomain = append(reviewDomain, *review.ToDomain())
	}

	return reviewDomain, nil
}
