package mentee_progresses

import (
	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
	"gorm.io/gorm"
)

type menteeProgressRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) menteeProgresses.Repository {
	return menteeProgressRepository{
		conn: conn,
	}
}

func (m menteeProgressRepository) Add(menteeProgressDomain *menteeProgresses.Domain) error {
	rec := FromDomain(menteeProgressDomain)

	err := m.conn.Model(&MenteeProgress{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (m menteeProgressRepository) FindByMaterial(menteeId string, materialId string) (*menteeProgresses.Domain, error) {
	rec := MenteeProgress{}

	err := m.conn.Model(&MenteeProgress{}).Where("mentee_progresses.mentee_id = ? AND mentee_progresses.material_id = ?", menteeId, materialId).
		First(&rec).Error

	if err != nil {
		return nil, err
	}

	return rec.ToDomain(), nil
}

func (m menteeProgressRepository) FindByMentee(menteeId string, courseId string) ([]menteeProgresses.Domain, error) {
	var rec []MenteeProgress

	err := m.conn.Model(&MenteeProgress{}).Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var progresses []menteeProgresses.Domain

	for _, progress := range rec {
		progresses = append(progresses, *progress.ToDomain())
	}

	return progresses, nil
}

func (m menteeProgressRepository) Count(menteeId string, title string, status string) ([]int64, error) {
	rec := []int64{}

	err := m.conn.Model(&MenteeProgress{}).Select("COUNT(DISTINCT mentee_progresses.material_id)").
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
	err := m.conn.Model(&MenteeProgress{}).Unscoped().
		Where("mentee_id = ? AND course_id = ?", menteeId, courseId).Delete(&MenteeProgress{}).Error

	if err != nil {
		return err
	}

	return nil
}
