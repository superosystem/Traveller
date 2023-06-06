package domain

import "time"

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryRepository interface {
	Create(categoryDomain *Category) error
	FindAll() (*[]Category, error)
	FindById(id string) (*Category, error)
	Update(id string, categoryDomain *Category) error
}

type CategoryUseCase interface {
	Create(categoryDomain *Category) error
	FindAll() (*[]Category, error)
	FindById(id string) (*Category, error)
	Update(id string, categoryDomain *Category) error
}
