package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type assignmentRepository struct {
	conn *gorm.DB
}

func NewAssignmentRepository(conn *gorm.DB) domain.AssignmentRepository {
	return assignmentRepository{
		conn: conn,
	}
}

func (ar assignmentRepository) Create(assignmentDomain *domain.Assignment) error {
	rec := entities.FromAssignmentDomain(assignmentDomain)

	err := ar.conn.Model(&entities.Assignment{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (ar assignmentRepository) FindById(assignmentId string) (*domain.Assignment, error) {
	rec := entities.Assignment{}

	err := ar.conn.Model(&entities.Assignment{}).Where("id = ?", assignmentId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrAssignmentNotFound
		}

		return nil, err
	}

	return rec.ToAssignmentDomain(), nil
}

func (ar assignmentRepository) FindByCourseId(courseId string) (*domain.Assignment, error) {
	rec := entities.Assignment{}

	err := ar.conn.Model(&entities.Assignment{}).Where("course_id = ?", courseId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCourseNotFound
		}

		return nil, err
	}

	return rec.ToAssignmentDomain(), nil
}

func (ar assignmentRepository) FindByCourses(courseIds []string) (*[]domain.Assignment, error) {
	rec := []entities.Assignment{}

	err := ar.conn.Model(&entities.Assignment{}).Where("course_id IN ?", courseIds).Find(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCourseNotFound
		}

		return nil, err
	}

	var assignmentDomain []domain.Assignment

	for _, assignment := range rec {
		assignmentDomain = append(assignmentDomain, *assignment.ToAssignmentDomain())
	}

	return &assignmentDomain, nil
}

func (ar assignmentRepository) Update(assignmentId string, assignmentDomain *domain.Assignment) error {
	rec := entities.FromAssignmentDomain(assignmentDomain)

	err := ar.conn.Model(&entities.Assignment{}).Where("id = ?", assignmentId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (ar assignmentRepository) Delete(assignmentId string) error {
	err := ar.conn.Model(&entities.Assignment{}).Where("id = ?", assignmentId).Delete(&entities.Assignment{}).Error

	if err != nil {
		return err
	}

	return nil
}
