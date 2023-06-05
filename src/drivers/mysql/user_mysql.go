package repository

import (
	"errors"

	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"gorm.io/gorm"
)

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) domain.UserRepository {
	return userRepository{
		conn: conn,
	}
}

func (ur userRepository) Create(userDomain *domain.User) error {
	rec := entities.FromUserDomain(userDomain)

	err := ur.conn.Model(&entities.User{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (ur userRepository) FindAll() (*[]domain.User, error) {
	panic("implement me")
}

func (ur userRepository) FindByEmail(email string) (*domain.User, error) {
	rec := entities.User{}

	err := ur.conn.Model(&entities.User{}).Where("email = ?", email).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrUserNotFound
		}

		return nil, err
	}

	return rec.ToUserDomain(), nil
}

func (ur userRepository) FindById(id string) (*domain.User, error) {
	rec := entities.User{}

	err := ur.conn.Model(&entities.User{}).Where("id = ?", id).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrUserNotFound
		}

		return nil, err
	}

	return rec.ToUserDomain(), nil
}

func (ur userRepository) Update(id string, userDomain *domain.User) error {
	rec := entities.FromUserDomain(userDomain)

	err := ur.conn.Model(&entities.User{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}
