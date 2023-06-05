package usecase

import (
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type categoryUsecase struct {
	categoryRepository domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepository domain.CategoryRepository) domain.CategoryUsecase {
	return categoryUsecase{
		categoryRepository: categoryRepository,
	}
}

func (c categoryUsecase) Create(categoryDomain *domain.Category) error {
	id := uuid.NewString()

	category := domain.Category{
		ID:   id,
		Name: categoryDomain.Name,
	}

	err := c.categoryRepository.Create(&category)

	if err != nil {
		return err
	}

	return nil
}

func (c categoryUsecase) FindAll() (*[]domain.Category, error) {
	categories, err := c.categoryRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c categoryUsecase) FindById(id string) (*domain.Category, error) {
	category, err := c.categoryRepository.FindById(id)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c categoryUsecase) Update(id string, categoryDomain *domain.Category) error {
	if _, err := c.categoryRepository.FindById(id); err != nil {
		return err
	}

	err := c.categoryRepository.Update(id, categoryDomain)

	if err != nil {
		return err
	}

	return nil
}

// func (c categoryUsecase) Delete(id string) error {
// 	//TODO implement me
// 	panic("implement me")
// }
