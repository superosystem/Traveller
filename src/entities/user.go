package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:200" json:"id"`
	Email     string    `gorm:"size:255" json:"email"`
	Password  string    `gorm:"size:255" json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromUserDomain(domain *domain.User) *User {
	return &User{
		ID:        domain.ID,
		Email:     domain.Email,
		Password:  domain.Password,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func (rec *User) ToUserDomain() *domain.User {
	return &domain.User{
		ID:        rec.ID,
		Email:     rec.Email,
		Password:  rec.Password,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}
}
