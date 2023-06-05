package domain

import "time"

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryRepository interface {
	// Create repository create categories
	Create(categoryDomain *Category) error

	// FindAll repository find all categories
	FindAll() (*[]Category, error)

	// FindById repository find by id categories
	FindById(id string) (*Category, error)

	// Update repository update categories
	Update(id string, categoryDomain *Category) error

	// Delete repository delete categories
	// Delete(id string) error
}

type CategoryUsecase interface {
	// Create usecase create categories
	Create(categoryDomain *Category) error

	// FindAll usecase find all categories
	FindAll() (*[]Category, error)

	// FindById usecase find categories by id
	FindById(id string) (*Category, error)

	// Update usecase update categories
	Update(id string, categoryDomain *Category) error

	// Delete usecase delete categories
	// Delete(id string) error
}
