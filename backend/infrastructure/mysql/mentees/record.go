package mentees

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/users"
)

type Mentee struct {
	ID             string     `gorm:"primaryKey;size:200" json:"id"`
	UserId         string     `gorm:"size:200" json:"user_id"`
	Fullname       string     `gorm:"size:255" json:"fullname"`
	Phone          string     `gorm:"size:15" json:"phone"`
	Role           string     `gorm:"size:50" json:"role"`
	BirthDate      string     `json:"birth_date"`
	ProfilePicture string     `json:"profile_picture"`
	User           users.User `json:"user"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func FromDomain(mentee *mentees.Domain) *Mentee {
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

func (rec *Mentee) ToDomain() *mentees.Domain {
	return &mentees.Domain{
		ID:             rec.ID,
		UserId:         rec.UserId,
		Fullname:       rec.Fullname,
		Phone:          rec.Phone,
		Role:           rec.Role,
		BirthDate:      rec.BirthDate,
		ProfilePicture: rec.ProfilePicture,
		User:           *rec.User.ToDomain(),
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}
