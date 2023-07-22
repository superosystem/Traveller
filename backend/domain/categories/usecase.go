package categories

import (
	"github.com/google/uuid"
)

type categoryUsecase struct {
	categoryRepository Repository
}

func NewCategoryUsecase(categoryRepository Repository) Usecase {
	return categoryUsecase{
		categoryRepository: categoryRepository,
	}
}

func (c categoryUsecase) Create(categoryDomain *Domain) error {
	id := uuid.NewString()

	category := Domain{
		ID:   id,
		Name: categoryDomain.Name,
	}

	err := c.categoryRepository.Create(&category)

	if err != nil {
		return err
	}

	return nil
}

func (c categoryUsecase) FindAll() (*[]Domain, error) {
	categories, err := c.categoryRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c categoryUsecase) FindById(id string) (*Domain, error) {
	category, err := c.categoryRepository.FindById(id)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c categoryUsecase) Update(id string, categoryDomain *Domain) error {
	if _, err := c.categoryRepository.FindById(id); err != nil {
		return err
	}

	err := c.categoryRepository.Update(id, categoryDomain)

	if err != nil {
		return err
	}

	return nil
}
