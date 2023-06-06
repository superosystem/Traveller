package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type moduleRepository struct {
	conn *gorm.DB
}

func NewModuleRepository(conn *gorm.DB) domain.ModuleRepository {
	return moduleRepository{
		conn: conn,
	}
}

func (mr moduleRepository) Create(moduleDomain *domain.Module) error {
	rec := entities.FromModuleDomain(moduleDomain)

	err := mr.conn.Model(&entities.Module{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr moduleRepository) FindById(moduleId string) (*domain.Module, error) {
	rec := entities.Module{}

	err := mr.conn.Model(&entities.Module{}).Where("id = ?", moduleId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrModuleNotFound
		}

		return nil, err
	}

	return rec.ToModuleDomain(), nil
}

func (mr moduleRepository) FindByCourse(courseId string) ([]domain.Module, error) {
	rec := []entities.Module{}

	err := mr.conn.Model(&entities.Module{}).Where("course_id = ?", courseId).
		Order("created_at ASC").
		Find(&rec).Error

	if err != nil {
		return nil, helper.ErrModuleNotFound
	}

	modulesDomain := []domain.Module{}

	for _, module := range rec {
		modulesDomain = append(modulesDomain, *module.ToModuleDomain())
	}

	return modulesDomain, nil
}

func (mr moduleRepository) Update(moduleId string, moduleDomain *domain.Module) error {
	rec := entities.FromModuleDomain(moduleDomain)

	err := mr.conn.Model(&entities.Module{}).Where("id = ?", moduleId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr moduleRepository) Delete(moduleId string) error {
	err := mr.conn.Model(&entities.Module{}).Where("id = ?", moduleId).Delete(&entities.Module{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr moduleRepository) Deletes(courseId string) error {
	err := mr.conn.Model(&entities.Module{}).Where("course_id = ?", courseId).Delete(&entities.Module{}).Error

	if err != nil {
		return err
	}

	return nil
}
