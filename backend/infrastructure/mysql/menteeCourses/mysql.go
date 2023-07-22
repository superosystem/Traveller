package mentee_courses

import (
	"errors"

	menteeCourses "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type menteeCourseRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) menteeCourses.Repository {
	return menteeCourseRepository{
		conn: conn,
	}
}

func (m menteeCourseRepository) Enroll(menteeCourseDomain *menteeCourses.Domain) error {
	rec := FromDomain(menteeCourseDomain)

	err := m.conn.Model(&MenteeCourse{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (m menteeCourseRepository) FindCoursesByMentee(menteeId string, title string, status string) (*[]menteeCourses.Domain, error) {
	var rec []MenteeCourse

	err := m.conn.Model(&MenteeCourse{}).Preload("Course.Mentor").
		Joins("INNER JOIN courses ON courses.id = mentee_courses.course_id").
		Where("mentee_courses.mentee_id = ? AND courses.title LIKE ? AND mentee_courses.status LIKE ? AND courses.deleted_at IS NULL", menteeId, "%"+title+"%", "%"+status+"%").
		Order("mentee_courses.course_id ASC").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var courses []menteeCourses.Domain

	for _, course := range rec {
		courses = append(courses, *course.ToDomain())
	}

	return &courses, nil
}

func (m menteeCourseRepository) CheckEnrollment(menteeId string, courseId string) (*menteeCourses.Domain, error) {
	rec := MenteeCourse{}

	err := m.conn.Model(&MenteeCourse{}).
		Where("mentee_courses.mentee_id = ? AND mentee_courses.course_id = ?", menteeId, courseId).
		First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrRecordNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (m menteeCourseRepository) Update(menteeId string, courseId string, menteeCourseDomain *menteeCourses.Domain) error {
	rec := FromDomain(menteeCourseDomain)

	err := m.conn.Model(&MenteeCourse{}).Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (m menteeCourseRepository) DeleteEnrolledCourse(menteeId string, courseId string) error {
	err := m.conn.Model(&MenteeCourse{}).Unscoped().
		Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Delete(&MenteeCourse{}).Error

	if err != nil {
		return err
	}

	return nil
}
