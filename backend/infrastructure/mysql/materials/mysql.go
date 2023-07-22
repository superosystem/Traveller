package materials

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type materialRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) materials.Repository {
	return materialRepository{
		conn: conn,
	}
}

func (mr materialRepository) Create(materialDomain *materials.Domain) error {
	rec := FromDomain(materialDomain)

	err := mr.conn.Model(&Material{}).Create(rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr materialRepository) FindById(materialId string) (*materials.Domain, error) {
	rec := Material{}

	err := mr.conn.Model(&Material{}).Preload("Module").
		Joins("INNER JOIN modules ON modules.id = materials.module_id").
		Where("materials.id = ?", materialId).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrMaterialAssetNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (mr materialRepository) FindByModule(moduleIds []string) ([]materials.Domain, error) {
	rec := []Material{}

	err := mr.conn.Model(&Material{}).Preload("Module").
		Joins("INNER JOIN modules ON modules.id = materials.module_id").
		Where("materials.module_id IN ?", moduleIds).
		Order("materials.created_at ASC").
		Find(&rec).Error

	if err != nil {
		return nil, err
	}

	materialsDomain := []materials.Domain{}

	for _, material := range rec {
		materialsDomain = append(materialsDomain, *material.ToDomain())
	}

	return materialsDomain, nil
}

func (mr materialRepository) CountByCourse(courseIds []string) ([]int64, error) {
	rec := []int64{}

	err := mr.conn.Model(&Material{}).Select("COUNT(materials.id)").
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

func (mr materialRepository) Update(materialId string, materialDomain *materials.Domain) error {
	rec := FromDomain(materialDomain)

	err := mr.conn.Model(&Material{}).Where("id = ?", materialId).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr materialRepository) Delete(materialId string) error {
	err := mr.conn.Model(&Material{}).Where("id = ?", materialId).Delete(&Material{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (mr materialRepository) Deletes(moduleId string) error {
	err := mr.conn.Model(&Material{}).Where("module_id = ?", moduleId).Delete(&Material{}).Error

	if err != nil {
		return err
	}

	return nil
}
