package categories

import (
	"errors"

	"github.com/superosystem/TrainingSystem/backend/domain/categories"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"gorm.io/gorm"
)

type categoryRepository struct {
	conn *gorm.DB
}

func NewSQLRepository(conn *gorm.DB) categories.Repository {
	return categoryRepository{
		conn: conn,
	}
}

func (cr categoryRepository) Create(categoryDomain *categories.Domain) error {
	rec := FromDomain(categoryDomain)

	err := cr.conn.Model(&Category{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return err
}

func (cr categoryRepository) FindAll() (*[]categories.Domain, error) {
	var rec []Category

	err := cr.conn.Model(&Category{}).Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var categories []categories.Domain

	for _, category := range rec {
		categories = append(categories, *category.ToDomain())
	}

	return &categories, nil
}

func (cr categoryRepository) FindById(id string) (*categories.Domain, error) {
	rec := Category{}

	err := cr.conn.Model(&Category{}).Where("id = ?", id).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCategoryNotFound
		}

		return nil, err
	}

	return rec.ToDomain(), nil
}

func (cr categoryRepository) Update(id string, categoryDomain *categories.Domain) error {
	rec := FromDomain(categoryDomain)

	err := cr.conn.Model(&Category{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}

// NOTE: optional
// func (cr categoryRepository) Delete(id string) error {
// 	panic("implement me")
// }
