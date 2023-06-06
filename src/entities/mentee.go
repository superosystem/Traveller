package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"time"
)

type Mentee struct {
	ID             string    `gorm:"primaryKey;size:200" json:"id"`
	UserId         string    `gorm:"size:200" json:"user_id"`
	Fullname       string    `gorm:"size:255" json:"fullname"`
	Phone          string    `gorm:"size:15" json:"phone"`
	Role           string    `gorm:"size:50" json:"role"`
	BirthDate      string    `json:"birth_date"`
	ProfilePicture string    `json:"profile_picture"`
	User           User      `json:"user"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (rec *Mentee) ToMenteeDomain() *domain.Mentee {
	return &domain.Mentee{
		ID:             rec.ID,
		UserId:         rec.UserId,
		Fullname:       rec.Fullname,
		Phone:          rec.Phone,
		Role:           rec.Role,
		BirthDate:      rec.BirthDate,
		ProfilePicture: rec.ProfilePicture,
		User:           *rec.User.ToUserDomain(),
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}

func FromMenteeDomain(mentee *domain.Mentee) *Mentee {
	return &Mentee{
		ID:             mentee.ID,
		UserId:         mentee.UserId,
		Fullname:       mentee.Fullname,
		Phone:          mentee.Phone,
		Role:           mentee.Role,
		BirthDate:      mentee.BirthDate,
		ProfilePicture: mentee.ProfilePicture,
		CreatedAt:      mentee.CreatedAt,
		UpdatedAt:      mentee.UpdatedAt,
	}
}
