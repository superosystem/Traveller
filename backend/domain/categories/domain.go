package categories

import "time"

type Domain struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	// Create repository create categories
	Create(categoryDomain *Domain) error

	// FindAll repository find all categories
	FindAll() (*[]Domain, error)

	// FindById repository find by id categories
	FindById(id string) (*Domain, error)

	// Update repository update categories
	Update(id string, categoryDomain *Domain) error
}

type Usecase interface {
	// Create usecase create categories
	Create(categoryDomain *Domain) error

	// FindAll usecase find all categories
	FindAll() (*[]Domain, error)

	// FindById usecase find categories by id
	FindById(id string) (*Domain, error)

	// Update usecase update categories
	Update(id string, categoryDomain *Domain) error
}
