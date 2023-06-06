package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type courseRepository struct {
	conn *gorm.DB
}

func NewCourseRepository(conn *gorm.DB) domain.CourseRepository {
	return courseRepository{
		conn: conn,
	}
}

func (cr courseRepository) Create(courseDomain *domain.Course) error {
	rec := entities.FromCourseDomain(courseDomain)

	err := cr.conn.Model(&entities.Course{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (cr courseRepository) FindAll(keyword string) (*[]domain.Course, error) {
	var rec []entities.CourseWithRating

	err := cr.conn.Model(&entities.Course{}).Preload("Category").Preload("Mentor").
		Select("COUNT(reviews.course_id) AS total_reviews, AVG(reviews.rating) as rating, courses.*, categories.id AS category_id, categories.name, mentors.id AS mentor_id, mentors.fullname").
		Joins("LEFT JOIN categories ON courses.category_id = categories.id").
		Joins("LEFT JOIN mentors ON courses.mentor_id = mentors.id").
		Joins("LEFT JOIN reviews ON courses.id = reviews.course_id").
		Where("courses.title LIKE ? OR categories.name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Group("courses.id").Order("courses.created_at").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var coursesDomain []domain.Course

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return &coursesDomain, nil
}

func (cr courseRepository) FindById(id string) (*domain.Course, error) {
	rec := entities.CourseWithRating{}

	err := cr.conn.Model(&entities.Course{}).Preload("Category").Preload("Mentor").
		Select("COUNT(reviews.course_id) AS total_reviews, AVG(reviews.rating) as rating, courses.*, categories.id AS category_id, categories.name, mentors.id AS mentor_id, mentors.fullname").
		Joins("LEFT JOIN categories ON courses.category_id = categories.id").
		Joins("LEFT JOIN mentors ON courses.mentor_id = mentors.id").
		Joins("LEFT JOIN reviews ON courses.id = reviews.course_id").
		Where("courses.id = ?", id).Group("courses.id").
		First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCourseNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (cr courseRepository) FindByCategory(categoryId string) (*[]domain.Course, error) {
	var rec []entities.CourseWithRating

	err := cr.conn.Model(&entities.Course{}).Preload("Category").Preload("Mentor").
		Select("COUNT(reviews.course_id) AS total_reviews, AVG(reviews.rating) as rating, courses.*, categories.id AS category_id, categories.name, mentors.id AS mentor_id, mentors.fullname").
		Joins("LEFT JOIN categories ON courses.category_id = categories.id").
		Joins("LEFT JOIN mentors ON courses.mentor_id = mentors.id").
		Joins("LEFT JOIN reviews ON courses.id = reviews.course_id").
		Where("courses.category_id = ?", categoryId).
		Group("courses.id").Order("courses.created_at").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var coursesDomain []domain.Course

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return &coursesDomain, nil
}

func (cr courseRepository) FindByMentor(mentorId string) (*[]domain.Course, error) {
	var rec []entities.CourseWithRating

	err := cr.conn.Model(&entities.Course{}).Preload("Category").Preload("Mentor").
		Select("COUNT(reviews.course_id) AS total_reviews, AVG(reviews.rating) as rating, courses.*, categories.id AS category_id, categories.name, mentors.id AS mentor_id, mentors.fullname").
		Joins("LEFT JOIN categories ON courses.category_id = categories.id").
		Joins("LEFT JOIN mentors ON courses.mentor_id = mentors.id").
		Joins("LEFT JOIN reviews ON courses.id = reviews.course_id").
		Where("courses.mentor_id = ?", mentorId).
		Group("courses.id").Order("courses.created_at").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var coursesDomain []domain.Course

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return &coursesDomain, nil
}

func (cr courseRepository) FindByPopular() ([]domain.Course, error) {
	var rec []entities.CourseWithRating

	err := cr.conn.Model(&entities.Course{}).Preload("Category").Preload("Mentor").
		Select("COUNT(reviews.course_id) AS total_reviews, AVG(reviews.rating) as rating, courses.*, categories.id AS category_id, categories.name, mentors.id AS mentor_id, mentors.fullname").
		Joins("LEFT JOIN categories ON courses.category_id = categories.id").
		Joins("LEFT JOIN mentors ON courses.mentor_id = mentors.id").
		Joins("LEFT JOIN reviews ON courses.id = reviews.course_id").
		Group("courses.id").Order("rating DESC").Limit(15).
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	coursesDomain := []domain.Course{}

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return coursesDomain, nil
}

func (cr courseRepository) Update(id string, courseDomain *domain.Course) error {
	rec := entities.FromCourseDomain(courseDomain)

	err := cr.conn.Model(&entities.Course{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (cr courseRepository) Delete(id string) error {
	err := cr.conn.Model(&entities.Course{}).Where("id = ?", id).Delete(&entities.Course{}).Error

	if err != nil {
		return err
	}

	return nil
}
