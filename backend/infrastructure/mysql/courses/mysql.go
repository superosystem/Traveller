package courses

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type courseRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) courses.Repository {
	return courseRepository{
		conn: conn,
	}
}

func (cr courseRepository) Create(courseDomain *courses.Domain) error {
	rec := FromDomain(courseDomain)

	err := cr.conn.Model(&Course{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (cr courseRepository) FindAll(keyword string) (*[]courses.Domain, error) {
	var rec []CourseWithRating

	err := cr.conn.Model(&Course{}).Preload("Category").Preload("Mentor").
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

	var coursesDomain []courses.Domain

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return &coursesDomain, nil
}

func (cr courseRepository) FindById(id string) (*courses.Domain, error) {
	rec := CourseWithRating{}

	err := cr.conn.Model(&Course{}).Preload("Category").Preload("Mentor").
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

func (cr courseRepository) FindByCategory(categoryId string) (*[]courses.Domain, error) {
	var rec []CourseWithRating

	err := cr.conn.Model(&Course{}).Preload("Category").Preload("Mentor").
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

	var coursesDomain []courses.Domain

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return &coursesDomain, nil
}

func (cr courseRepository) FindByMentor(mentorId string) (*[]courses.Domain, error) {
	var rec []CourseWithRating

	err := cr.conn.Model(&Course{}).Preload("Category").Preload("Mentor").
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

	var coursesDomain []courses.Domain

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return &coursesDomain, nil
}

func (cr courseRepository) FindByPopular() ([]courses.Domain, error) {
	var rec []CourseWithRating

	err := cr.conn.Model(&Course{}).Preload("Category").Preload("Mentor").
		Select("COUNT(reviews.course_id) AS total_reviews, AVG(reviews.rating) as rating, courses.*, categories.id AS category_id, categories.name, mentors.id AS mentor_id, mentors.fullname").
		Joins("LEFT JOIN categories ON courses.category_id = categories.id").
		Joins("LEFT JOIN mentors ON courses.mentor_id = mentors.id").
		Joins("LEFT JOIN reviews ON courses.id = reviews.course_id").
		Group("courses.id").Order("rating DESC").Limit(15).
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	coursesDomain := []courses.Domain{}

	for _, course := range rec {
		coursesDomain = append(coursesDomain, *course.ToDomain())
	}

	return coursesDomain, nil
}

func (cr courseRepository) Update(id string, courseDomain *courses.Domain) error {
	rec := FromDomain(courseDomain)

	err := cr.conn.Model(&Course{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (cr courseRepository) Delete(id string) error {
	err := cr.conn.Model(&Course{}).Where("id = ?", id).Delete(&Course{}).Error

	if err != nil {
		return err
	}

	return nil
}
