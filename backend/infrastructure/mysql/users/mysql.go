package users

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/users"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type userRepository struct {
	conn *gorm.DB
}

func NewMySQLRepository(conn *gorm.DB) users.Repository {
	return userRepository{
		conn: conn,
	}
}

func (ur userRepository) Create(userDomain *users.Domain) error {
	rec := FromDomain(userDomain)

	err := ur.conn.Model(&User{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (ur userRepository) FindAll() (*[]users.Domain, error) {
	panic("implement me")
}

func (ur userRepository) FindByEmail(email string) (*users.Domain, error) {
	rec := User{}

	err := ur.conn.Model(&User{}).Where("email = ?", email).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrUserNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (ur userRepository) FindById(id string) (*users.Domain, error) {
	rec := User{}

	err := ur.conn.Model(&User{}).Where("id = ?", id).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrUserNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (ur userRepository) Update(id string, userDomain *users.Domain) error {
	rec := FromDomain(userDomain)

	err := ur.conn.Model(&User{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}
