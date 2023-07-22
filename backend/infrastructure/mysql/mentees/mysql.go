package mentees

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type menteeRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) mentees.Repository {
	return menteeRepository{
		conn: conn,
	}
}

func (mr menteeRepository) Create(menteeDomain *mentees.Domain) error {
	rec := FromDomain(menteeDomain)

	err := mr.conn.Model(&Mentee{}).Omit("birth_date", "address", "profile_picture").Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr menteeRepository) FindAll() (*[]mentees.Domain, error) {
	var rec []Mentee

	err := mr.conn.Model(&Mentee{}).Preload("User").Find(&rec).Error

	if err != nil {
		return nil, err
	}

	menteeDomain := []mentees.Domain{}

	for _, mentee := range rec {
		menteeDomain = append(menteeDomain, *mentee.ToDomain())
	}

	return &menteeDomain, nil
}

func (mr menteeRepository) FindById(id string) (*mentees.Domain, error) {
	rec := Mentee{}

	err := mr.conn.Model(&Mentee{}).Where("id = ?", id).Preload("User").First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMenteeNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (mr menteeRepository) FindByIdUser(userId string) (*mentees.Domain, error) {
	rec := Mentee{}

	err := mr.conn.Model(&Mentee{}).Where("user_id = ?", userId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMenteeNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (mr menteeRepository) FindByCourse(courseId string, limit int, offset int) (*[]mentees.Domain, int, error) {
	var totalRows int64

	_ = mr.conn.Model(&Mentee{}).
		Joins("LEFT JOIN users ON users.id = mentees.user_id").
		Joins("LEFT JOIN mentee_courses ON mentees.id = mentee_courses.mentee_id").
		Where("mentee_courses.course_id = ?", courseId).
		Order("mentees.fullname ASC").
		Count(&totalRows).Error

	var rec []Mentee

	err := mr.conn.Model(&Mentee{}).Preload("User").
		Select("mentees.id, mentees.user_id, mentees.fullname, mentees.phone, mentees.role, mentees.birth_date, mentees.profile_picture, users.email, mentee_courses.created_at, mentee_courses.updated_at").
		Joins("LEFT JOIN users ON users.id = mentees.user_id").
		Joins("LEFT JOIN mentee_courses ON mentees.id = mentee_courses.mentee_id").
		Where("mentee_courses.course_id = ?", courseId).
		Order("mentees.fullname ASC").Limit(limit).Offset(offset).
		Find(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, helper.ErrCourseNotFound
		}

		return nil, 0, err
	}

	var menteeDomain []mentees.Domain

	for _, mentee := range rec {
		menteeDomain = append(menteeDomain, *mentee.ToDomain())
	}

	return &menteeDomain, int(totalRows), nil
}

func (mr menteeRepository) CountByCourse(courseId string) (int64, error) {
	var total int64

	err := mr.conn.Model(&Mentee{}).
		Joins("LEFT JOIN users ON users.id = mentees.user_id").
		Joins("LEFT JOIN mentee_courses ON mentees.id = mentee_courses.mentee_id").
		Where("mentee_courses.course_id = ?", courseId).
		Order("mentees.fullname ASC").
		Count(&total).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, helper.ErrCourseNotFound
		}

		return 0, nil
	}

	return total, nil
}

func (mr menteeRepository) Update(id string, menteeDomain *mentees.Domain) error {
	rec := FromDomain(menteeDomain)

	err := mr.conn.Model(&Mentee{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}
