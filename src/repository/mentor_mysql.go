package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type mentorRepository struct {
	conn *gorm.DB
}

func NewMentorRepository(conn *gorm.DB) domain.MentorRepository {
	return mentorRepository{
		conn: conn,
	}
}

func (mr mentorRepository) Create(mentorDomain *domain.Mentor) error {
	rec := entities.FromMentorDomain(mentorDomain)

	err := mr.conn.Model(&entities.Mentor{}).Omit("jobs",
		"gender", "phone", "birth_place", "birth_date", "address", "profile_picture").Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr mentorRepository) FindAll() (*[]domain.Mentor, error) {
	var rec []entities.Mentor

	err := mr.conn.Model(&entities.Mentor{}).Preload("User").Find(&rec).Error

	if err != nil {
		return nil, err
	}

	mentorDomain := []domain.Mentor{}

	for _, mentor := range rec {
		mentorDomain = append(mentorDomain, *mentor.ToMentorDomain())
	}

	return &mentorDomain, nil
}

func (mr mentorRepository) FindById(id string) (*domain.Mentor, error) {
	rec := entities.Mentor{}

	err := mr.conn.Model(&entities.Mentor{}).Where("id = ?", id).Preload("User").First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMentorNotFound
		}

		return nil, err
	}

	return rec.ToMentorDomain(), nil
}

func (mr mentorRepository) FindByIdUser(userId string) (*domain.Mentor, error) {
	rec := &entities.Mentor{}

	err := mr.conn.Model(&entities.Mentor{}).Where("user_id = ?", userId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMentorNotFound
		}

		return nil, err
	}

	return rec.ToMentorDomain(), nil
}

func (mr mentorRepository) Update(id string, mentorDomain *domain.Mentor) error {
	rec := entities.FromMentorDomain(mentorDomain)

	err := mr.conn.Model(&entities.Mentor{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}
