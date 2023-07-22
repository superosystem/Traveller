package users

import "time"

type Domain struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	// Create add new user
	Create(userDomain *Domain) error

	// FindAll find all users
	FindAll() (*[]Domain, error)

	// FindByEmail find user by email
	FindByEmail(email string) (*Domain, error)

	// FindById find user by id
	FindById(id string) (*Domain, error)

	// Update edit data user
	Update(id string, userDomain *Domain) error
}

type Usecase interface{}
