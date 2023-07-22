package modules

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/modules"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type moduleRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) modules.Repository {
	return moduleRepository{
		conn: conn,
	}
}

func (mr moduleRepository) Create(moduleDomain *modules.Domain) error {
	rec := FromDomain(moduleDomain)

	err := mr.conn.Model(&Module{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr moduleRepository) FindById(moduleId string) (*modules.Domain, error) {
	rec := Module{}

	err := mr.conn.Model(&Module{}).Where("id = ?", moduleId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrModuleNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (mr moduleRepository) FindByCourse(courseId string) ([]modules.Domain, error) {
	rec := []Module{}

	err := mr.conn.Model(&Module{}).Where("course_id = ?", courseId).
		Order("created_at ASC").
		Find(&rec).Error

	if err != nil {
		return nil, helper.ErrModuleNotFound
	}

	modulesDomain := []modules.Domain{}

	for _, module := range rec {
		modulesDomain = append(modulesDomain, *module.ToDomain())
	}

	return modulesDomain, nil
}

func (mr moduleRepository) Update(moduleId string, moduleDomain *modules.Domain) error {
	rec := FromDomain(moduleDomain)

	err := mr.conn.Model(&Module{}).Where("id = ?", moduleId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr moduleRepository) Delete(moduleId string) error {
	err := mr.conn.Model(&Module{}).Where("id = ?", moduleId).Delete(&Module{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr moduleRepository) Deletes(courseId string) error {
	err := mr.conn.Model(&Module{}).Where("course_id = ?", courseId).Delete(&Module{}).Error

	if err != nil {
		return err
	}

	return nil
}
