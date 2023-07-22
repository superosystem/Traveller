package users

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/users"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:200" json:"id"`
	Email     string    `gorm:"size:255" json:"email"`
	Password  string    `gorm:"size:255" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromDomain(domain *users.Domain) *User {
	return &User{
		ID:        domain.ID,
		Email:     domain.Email,
		Password:  domain.Password,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (rec *User) ToDomain() *users.Domain {
	return &users.Domain{
		ID:        rec.ID,
		Email:     rec.Email,
		Password:  rec.Password,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}
