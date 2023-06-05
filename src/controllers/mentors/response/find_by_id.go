package response

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type FindMentorByID struct {
	ID             string    `json:"id,omitempty"`
	UserID         string    `json:"user_id"`
	Fullname       string    `json:"fullname"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Role           string    `json:"role"`
	Jobs           string    `json:"jobs"`
	Gender         string    `json:"gender"`
	BirthPlace     string    `json:"birth_place"`
	BirthDate      string    `json:"birth_date"`
	ProfilePicture string    `json:"profile_picture"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func FromDomainUser(user *domain.Mentor) *FindMentorByID {
	return &FindMentorByID{
		ID:             user.ID,
		UserID:         user.UserId,
		Fullname:       user.Fullname,
		Email:          user.Email,
		Phone:          user.Phone,
		Role:           user.Role,
		Jobs:           user.Jobs,
		Gender:         user.Gender,
		BirthPlace:     user.BirthPlace,
		BirthDate:      user.BirthDate.String(),
		ProfilePicture: user.ProfilePicture,
		Address:        user.Address,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
