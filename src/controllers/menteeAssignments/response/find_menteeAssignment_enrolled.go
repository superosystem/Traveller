package response

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type AssignmentMenteeEnrolled struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	PDF            string    `json:"pdf"`
	ProfilePicture string    `json:"profile_picture"`
	Grade          int       `json:"grade"`
	Completed      bool      `json:"completed"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func DetailAssignmentEnrolled(domain *domain.MenteeAssignment) AssignmentMenteeEnrolled {
	return AssignmentMenteeEnrolled{
		ID:             domain.ID,
		Name:           domain.Name,
		PDF:            domain.AssignmentURL,
		ProfilePicture: domain.ProfilePicture,
		Grade:          domain.Grade,
		Completed:      domain.Completed,
		CreatedAt:      domain.CreatedAt,
		UpdatedAt:      domain.UpdatedAt,
	}
}
