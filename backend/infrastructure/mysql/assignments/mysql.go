package assignments

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/assignments"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type assignmentRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) assignments.Repository {
	return assignmentRepository{
		conn: conn,
	}
}

func (ar assignmentRepository) Create(assignmentDomain *assignments.Domain) error {
	rec := FromDomain(assignmentDomain)

	err := ar.conn.Model(&Assignment{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (ar assignmentRepository) FindById(assignmentId string) (*assignments.Domain, error) {
	rec := Assignment{}

	err := ar.conn.Model(&Assignment{}).Where("id = ?", assignmentId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrAssignmentNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (ar assignmentRepository) FindByCourseId(courseId string) (*assignments.Domain, error) {
	rec := Assignment{}

	err := ar.conn.Model(&Assignment{}).Where("course_id = ?", courseId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCourseNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (ar assignmentRepository) FindByCourses(courseIds []string) (*[]assignments.Domain, error) {
	rec := []Assignment{}

	err := ar.conn.Model(&Assignment{}).Where("course_id IN ?", courseIds).Find(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCourseNotFound
		}

		return nil, err
	}

	var assignmentDomain []assignments.Domain

	for _, assignment := range rec {
		assignmentDomain = append(assignmentDomain, *assignment.ToDomain())
	}

	return &assignmentDomain, nil
}

func (ar assignmentRepository) Update(assignmentId string, assignmentDomain *assignments.Domain) error {
	rec := FromDomain(assignmentDomain)

	err := ar.conn.Model(&Assignment{}).Where("id = ?", assignmentId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (ar assignmentRepository) Delete(assignmentId string) error {
	err := ar.conn.Model(&Assignment{}).Where("id = ?", assignmentId).Delete(&Assignment{}).Error

	if err != nil {
		return err
	}

	return nil
}
