package mentors

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type mentorRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) mentors.Repository {
	return mentorRepository{
		conn: conn,
	}
}

func (mr mentorRepository) Create(mentorDomain *mentors.Domain) error {
	rec := FromDomain(mentorDomain)

	err := mr.conn.Model(&Mentor{}).Omit("jobs",
		"gender", "phone", "birth_place", "birth_date", "address", "profile_picture").Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr mentorRepository) FindAll() (*[]mentors.Domain, error) {
	var rec []Mentor

	err := mr.conn.Model(&Mentor{}).Preload("User").Find(&rec).Error

	if err != nil {
		return nil, err
	}

	mentorDomain := []mentors.Domain{}

	for _, mentor := range rec {
		mentorDomain = append(mentorDomain, *mentor.ToDomain())
	}

	return &mentorDomain, nil
}

func (mr mentorRepository) FindById(id string) (*mentors.Domain, error) {
	rec := Mentor{}

	err := mr.conn.Model(&Mentor{}).Where("id = ?", id).Preload("User").First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMentorNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (mr mentorRepository) FindByIdUser(userId string) (*mentors.Domain, error) {
	rec := &Mentor{}

	err := mr.conn.Model(&Mentor{}).Where("user_id = ?", userId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMentorNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (mr mentorRepository) Update(id string, mentorDomain *mentors.Domain) error {
	rec := FromDomain(mentorDomain)

	err := mr.conn.Model(&Mentor{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}
