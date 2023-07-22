package mentee_assignments

import (
	"time"

	menteeAssignments "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/assignments"
	"github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/mentees"
)

type MenteeAssignment struct {
	ID            string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId      string `json:"mentee_id" gorm:"size:200"`
	AssignmentId  string `json:"assignment_id" gorm:"size:200"`
	AssignmentURL string `json:"assignment_url"`
	Grade         int    `json:"grade"`
	Mentee        mentees.Mentee
	Assignment    assignments.Assignment
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (rec *MenteeAssignment) ToDomain() *menteeAssignments.Domain {
	return &menteeAssignments.Domain{
		ID:             rec.ID,
		MenteeId:       rec.MenteeId,
		AssignmentId:   rec.AssignmentId,
		Name:           rec.Mentee.Fullname,
		ProfilePicture: rec.Mentee.ProfilePicture,
		AssignmentURL:  rec.AssignmentURL,
		Grade:          rec.Grade,
		Mentee:         *rec.Mentee.ToDomain(),
		Assignment:     *rec.Assignment.ToDomain(),
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}

func FromDomain(menteeAssignmentDomain *menteeAssignments.Domain) *MenteeAssignment {
	return &MenteeAssignment{
		ID:            menteeAssignmentDomain.ID,
		MenteeId:      menteeAssignmentDomain.MenteeId,
		AssignmentId:  menteeAssignmentDomain.AssignmentId,
		AssignmentURL: menteeAssignmentDomain.AssignmentURL,
		Grade:         menteeAssignmentDomain.Grade,
		CreatedAt:     menteeAssignmentDomain.CreatedAt,
		UpdatedAt:     menteeAssignmentDomain.UpdatedAt,
	}
}
