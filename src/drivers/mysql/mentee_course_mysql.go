package repository

import (
	"errors"

	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"gorm.io/gorm"
)

type menteeCourseRepository struct {
	conn *gorm.DB
}

func NewMenteeCourseRepository(conn *gorm.DB) domain.MenteeCourseRepository {
	return menteeCourseRepository{
		conn: conn,
	}
}

func (m menteeCourseRepository) Enroll(menteeCourseDomain *domain.MenteeCourse) error {
	rec := entities.FromMenteeCourseDomain(menteeCourseDomain)

	err := m.conn.Model(&entities.MenteeCourse{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (m menteeCourseRepository) FindCoursesByMentee(menteeId string, title string, status string) (*[]domain.MenteeCourse, error) {
	var rec []entities.MenteeCourse

	err := m.conn.Model(&entities.MenteeCourse{}).Preload("Course.Mentor").
		Joins("INNER JOIN courses ON courses.id = mentee_courses.course_id").
		Where("mentee_courses.mentee_id = ? AND courses.title LIKE ? AND mentee_courses.status LIKE ? AND courses.deleted_at IS NULL", menteeId, "%"+title+"%", "%"+status+"%").
		Order("mentee_courses.course_id ASC").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var courses []domain.MenteeCourse

	for _, course := range rec {
		courses = append(courses, *course.ToMenteeCourseDomain())
	}

	return &courses, nil
}

func (m menteeCourseRepository) CheckEnrollment(menteeId string, courseId string) (*domain.MenteeCourse, error) {
	rec := entities.MenteeCourse{}

	err := m.conn.Model(&entities.MenteeCourse{}).
		Where("mentee_courses.mentee_id = ? AND mentee_courses.course_id = ?", menteeId, courseId).
		First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrRecordNotFound
		}

		return nil, err
	}

	return rec.ToMenteeCourseDomain(), nil
}

func (m menteeCourseRepository) Update(menteeId string, courseId string, menteeCourseDomain *domain.MenteeCourse) error {
	rec := entities.FromMenteeCourseDomain(menteeCourseDomain)

	err := m.conn.Model(&entities.MenteeCourse{}).Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (m menteeCourseRepository) DeleteEnrolledCourse(menteeId string, courseId string) error {
	err := m.conn.Model(&entities.MenteeCourse{}).Unscoped().
		Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Delete(&entities.MenteeCourse{}).Error

	if err != nil {
		return err
	}

	return nil
}
