package response

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type FindByIdMentees struct {
	ID             string    `json:"id"`
	UserId         string    `json:"user_id"`
	Fullname       string    `json:"fullname"`
	Phone          string    `json:"phone"`
	Role           string    `json:"role"`
	BirthDate      string    `json:"birth_date"`
	ProfilePicture string    `json:"profile_picture"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func FromDomain(menteeDomain *domain.Mentee) *FindByIdMentees {
	return &FindByIdMentees{
		ID:             menteeDomain.ID,
		UserId:         menteeDomain.UserId,
		Fullname:       menteeDomain.Fullname,
		Phone:          menteeDomain.Phone,
		Role:           menteeDomain.Role,
		BirthDate:      menteeDomain.BirthDate,
		ProfilePicture: menteeDomain.ProfilePicture,
		Email:          menteeDomain.User.Email,
		CreatedAt:      menteeDomain.CreatedAt,
		UpdatedAt:      menteeDomain.UpdatedAt,
	}
}