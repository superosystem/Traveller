package repository

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type reviewRepository struct {
	conn *gorm.DB
}

func NewReviewRepository(conn *gorm.DB) domain.ReviewRepository {
	return reviewRepository{
		conn: conn,
	}
}

func (rr reviewRepository) Create(reviewDomain *domain.Review) error {
	rec := entities.FromReviewDomain(reviewDomain)

	err := rr.conn.Model(&entities.Review{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (rr reviewRepository) FindByCourse(courseId string) ([]domain.Review, error) {
	var rec []entities.Review

	err := rr.conn.Model(&entities.Review{}).Preload("Mentee").Preload("Course").
		Joins("INNER JOIN mentees ON mentees.id = reviews.mentee_id").
		Joins("INNER JOIN courses ON courses.id = reviews.course_id").
		Where("reviews.course_id = ?", courseId).
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var reviewDomain []domain.Review

	for _, review := range rec {
		reviewDomain = append(reviewDomain, *review.ToReviewDomain())
	}

	return reviewDomain, nil
}
