package repository

import (
	"errors"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/gorm"
)

type categoryRepository struct {
	conn *gorm.DB
}

func NewCategoryRepository(conn *gorm.DB) domain.CategoryRepository {
	return categoryRepository{
		conn: conn,
	}
}

func (cr categoryRepository) Create(categoryDomain *domain.Category) error {
	rec := entities.FromCategoryDomain(categoryDomain)

	err := cr.conn.Model(&entities.Category{}).Create(&rec).Error

	if err != nil {
		return err
	}

	return err
}

func (cr categoryRepository) FindAll() (*[]domain.Category, error) {
	var rec []entities.Category

	err := cr.conn.Model(&entities.Category{}).Find(&rec).Error

	if err != nil {
		return nil, err
	}

	var categories []domain.Category

	for _, category := range rec {
		categories = append(categories, *category.ToCategoryDomain())
	}

	return &categories, nil
}

func (cr categoryRepository) FindById(id string) (*domain.Category, error) {
	rec := entities.Category{}

	err := cr.conn.Model(&entities.Category{}).Where("id = ?", id).First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrCategoryNotFound
		}

		return nil, err
	}

	return rec.ToCategoryDomain(), nil
}

func (cr categoryRepository) Update(id string, categoryDomain *domain.Category) error {
	rec := entities.FromCategoryDomain(categoryDomain)

	err := cr.conn.Model(&entities.Category{}).Where("id = ?", id).Updates(&rec).Error

	if err != nil {
		return err
	}

	return nil
}
