package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type materialRepository struct {
	conn *gorm.DB
}

func NewMaterialRepository(conn *gorm.DB) domain.MaterialRepository {
	return materialRepository{
		conn: conn,
	}
}

func (mr materialRepository) Create(materialDomain *domain.Material) error {
	rec := entities.FromMaterialDomain(materialDomain)

	err := mr.conn.Model(&entities.Material{}).Create(rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr materialRepository) FindById(materialId string) (*domain.Material, error) {
	rec := entities.Material{}

	err := mr.conn.Model(&entities.Material{}).Preload("Module").
		Joins("INNER JOIN modules ON modules.id = materials.module_id").
		Where("materials.id = ?", materialId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMaterialAssetNotFound
		}

		return nil, err
	}

	return rec.ToMaterialDomain(), nil
}

func (mr materialRepository) FindByModule(moduleIds []string) ([]domain.Material, error) {
	rec := []entities.Material{}

	err := mr.conn.Model(&entities.Material{}).Preload("Module").
		Joins("INNER JOIN modules ON modules.id = materials.module_id").
		Where("materials.module_id IN ?", moduleIds).
		Order("materials.created_at ASC").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	materialsDomain := []domain.Material{}

	for _, material := range rec {
		materialsDomain = append(materialsDomain, *material.ToMaterialDomain())
	}

	return materialsDomain, nil
}

func (mr materialRepository) CountByCourse(courseIds []string) ([]int64, error) {
	rec := []int64{}

	err := mr.conn.Model(&entities.Material{}).Select("COUNT(materials.id)").
		Joins("INNER JOIN modules ON modules.id = materials.module_id").
		Where("modules.course_id IN ?", courseIds).
		Group("modules.course_id").
		Order("modules.course_id ASC").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (mr materialRepository) Update(materialId string, materialDomain *domain.Material) error {
	rec := entities.FromMaterialDomain(materialDomain)

	err := mr.conn.Model(&entities.Material{}).Where("id = ?", materialId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr materialRepository) Delete(materialId string) error {
	err := mr.conn.Model(&entities.Material{}).Where("id = ?", materialId).Delete(&entities.Material{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr materialRepository) Deletes(moduleId string) error {
	err := mr.conn.Model(&entities.Material{}).Where("module_id = ?", moduleId).Delete(&entities.Material{}).Error

	if err != nil {
		return err
	}

	return nil
}
