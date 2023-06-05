package domain

import "time"

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	// Create add new user
	Create(userDomain *User) error

	// FindAll find all users
	FindAll() (*[]User, error)

	// FindByEmail find user by email
	FindByEmail(email string) (*User, error)

	// FindById find user by id
	FindById(id string) (*User, error)

	// Update edit data user
	Update(id string, userDomain *User) error
}

type UserUsecase interface{}
