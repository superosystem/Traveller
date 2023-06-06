package repository

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type menteeProgressRepository struct {
	conn *gorm.DB
}

func NewMenteeProgressRepository(conn *gorm.DB) domain.MenteeProgressRepository {
	return menteeProgressRepository{
		conn: conn,
	}
}

func (m menteeProgressRepository) Add(menteeProgressDomain *domain.MenteeProgress) error {
	rec := entities.FromMenteeProgressDomain(menteeProgressDomain)

	err := m.conn.Model(&entities.MenteeProgress{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (m menteeProgressRepository) FindByMaterial(menteeId string, materialId string) (*domain.MenteeProgress, error) {
	rec := entities.MenteeProgress{}

	err := m.conn.Model(&entities.MenteeProgress{}).Where("mentee_progresses.mentee_id = ? AND mentee_progresses.material_id = ?", menteeId, materialId).
		First(&rec).Error

	if err != nil {
		return nil, err
	}

	return rec.ToMenteeProgressDomain(), nil
}

func (m menteeProgressRepository) FindByMentee(menteeId string, courseId string) ([]domain.MenteeProgress, error) {
	var rec []entities.MenteeProgress

	err := m.conn.Model(&entities.MenteeProgress{}).Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var progresses []domain.MenteeProgress

	for _, progress := range rec {
		progresses = append(progresses, *progress.ToMenteeProgressDomain())
	}

	return progresses, nil
}

func (m menteeProgressRepository) Count(menteeId string, title string, status string) ([]int64, error) {
	rec := []int64{}

	err := m.conn.Model(&entities.MenteeProgress{}).Select("COUNT(DISTINCT mentee_progresses.material_id)").
		Joins("LEFT JOIN courses ON courses.id = mentee_progresses.course_id").
		Joins("LEFT JOIN mentee_courses ON courses.id = mentee_courses.course_id").
		Where("mentee_progresses.mentee_id = ? AND courses.title LIKE ? AND mentee_courses.status LIKE ?", menteeId, "%"+title+"%", "%"+status+"%").
		Group("mentee_progresses.mentee_id").Group("mentee_progresses.course_id").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (m menteeProgressRepository) DeleteMenteeProgressesByCourse(menteeId string, courseId string) error {
	err := m.conn.Model(&entities.MenteeProgress{}).Unscoped().
		Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Delete(&entities.MenteeProgress{}).Error

	if err != nil {
		return err
	}

	return nil
}
