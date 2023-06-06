package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type assignmentMenteeRepository struct {
	conn *gorm.DB
}

func NewMenteeAssignmentRepository(conn *gorm.DB) domain.MenteeAssignmentRepository {
	return assignmentMenteeRepository{
		conn: conn,
	}
}

func (am assignmentMenteeRepository) Create(assignmentmenteeDomain *domain.MenteeAssignment) error {
	rec := entities.FromMenteeAssignmentDomain(assignmentmenteeDomain)

	err := am.conn.Model(&entities.MenteeAssignment{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (am assignmentMenteeRepository) FindById(assignmentMenteeId string) (*domain.MenteeAssignment, error) {
	rec := entities.MenteeAssignment{}

	err := am.conn.Model(&entities.MenteeAssignment{}).Where("id = ?", assignmentMenteeId).Preload("Mentee").First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, err
	}

	return rec.ToMenteeAssignmentDomain(), nil
}

func (am assignmentMenteeRepository) FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*domain.MenteeAssignment, error) {
	rec := entities.MenteeAssignment{}

	err := am.conn.Model(&entities.MenteeAssignment{}).Where("mentee_assignments.mentee_id = ? AND mentee_assignments.assignment_id  = ?", menteeId, assignmentId).Preload("Mentee").
		First(&rec).Error

	if err != nil {
		return nil, err
	}

	return rec.ToMenteeAssignmentDomain(), nil
}

func (am assignmentMenteeRepository) FindByMenteeId(menteeId string) ([]domain.MenteeAssignment, error) {
	rec := []entities.MenteeAssignment{}

	err := am.conn.Model(&entities.MenteeAssignment{}).Where("mentee_id = ?", menteeId).Preload("Mentee").Order("created_at DESC").Find(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, err
	}
	assignmentMenteeDomain := []domain.MenteeAssignment{}

	for _, assignment := range rec {
		assignmentMenteeDomain = append(assignmentMenteeDomain, *assignment.ToMenteeAssignmentDomain())
	}

	return assignmentMenteeDomain, nil

}

func (am assignmentMenteeRepository) FindByAssignmentId(assignmentId string, limit int, offset int) ([]domain.MenteeAssignment, int, error) {
	var totalRows int64

	_ = am.conn.Model(&entities.MenteeAssignment{}).
		Where("assignment_id = ?", assignmentId).Order("created_at DESC").
		Count(&totalRows).Error

	rec := []entities.MenteeAssignment{}

	err := am.conn.Model(&entities.MenteeAssignment{}).Preload("Mentee").
		Where("assignment_id = ?", assignmentId).Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, helper.ErrAssignmentNotFound
		}

		return nil, 0, err
	}

	assignmentDomain := []domain.MenteeAssignment{}

	for _, assignment := range rec {
		assignmentDomain = append(assignmentDomain, *assignment.ToMenteeAssignmentDomain())
	}

	return assignmentDomain, int(totalRows), nil
}

func (am assignmentMenteeRepository) FindByCourse(menteeId string, courseId string) (*domain.MenteeAssignment, error) {
	rec := entities.MenteeAssignment{}

	err := am.conn.Model(&entities.MenteeAssignment{}).Preload("Assignment").
		Joins("LEFT JOIN assignments ON assignments.id = mentee_assignments.assignment_id").
		Where("mentee_assignments.mentee_id = ? AND assignments.course_id = ?", menteeId, courseId).
		First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrAssignmentNotFound
		}

		return nil, err
	}

	return rec.ToMenteeAssignmentDomain(), nil
}

func (am assignmentMenteeRepository) FindByCourses(menteeId string, courseIds []string) (*[]domain.MenteeAssignment, error) {
	rec := []entities.MenteeAssignment{}

	err := am.conn.Model(&entities.MenteeAssignment{}).Preload("Assignment").
		Joins("LEFT JOIN assignments ON assignments.id = mentee_assignments.assignment_id").
		Where("mentee_assignments.mentee_id = ? AND assignments.course_id IN ?", menteeId, courseIds).
		First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrAssignmentNotFound
		}

		return nil, err
	}

	var menteeAssignmentDomain []domain.MenteeAssignment

	for _, assignment := range rec {
		menteeAssignmentDomain = append(menteeAssignmentDomain, *assignment.ToMenteeAssignmentDomain())
	}

	return &menteeAssignmentDomain, nil
}

func (am assignmentMenteeRepository) Update(assignmentMenteeId string, assignmentmenteeDomain *domain.MenteeAssignment) error {
	rec := entities.FromMenteeAssignmentDomain(assignmentmenteeDomain)

	err := am.conn.Model(&entities.MenteeAssignment{}).Where("id = ?", assignmentMenteeId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (am assignmentMenteeRepository) Delete(assignmentMenteeId string) error {
	err := am.conn.Model(&entities.MenteeAssignment{}).Unscoped().Where("id = ?", assignmentMenteeId).Delete(&entities.MenteeAssignment{}).Error

	if err != nil {
		return err
	}

	return nil
}
