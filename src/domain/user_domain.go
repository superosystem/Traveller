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
	Create(userDomain *User) error
	FindAll() (*[]User, error)
	FindByEmail(email string) (*User, error)
	FindById(id string) (*User, error)
	Update(id string, userDomain *User) error
}

type UserUseCase interface{}
