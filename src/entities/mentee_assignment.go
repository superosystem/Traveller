package entities

import (
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"time"
)

type MenteeAssignment struct {
	ID            string `json:"id" gorm:"primaryKey;size:200"`
	MenteeId      string `json:"mentee_id" gorm:"size:200"`
	AssignmentId  string `json:"assignment_id" gorm:"size:200"`
	AssignmentURL string `json:"assignment_url"`
	Grade         int    `json:"grade"`
	Mentee        Mentee
	Assignment    Assignment
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (rec *MenteeAssignment) ToMenteeAssignmentDomain() *domain.MenteeAssignment {
	return &domain.MenteeAssignment{
		ID:             rec.ID,
		MenteeId:       rec.MenteeId,
		AssignmentId:   rec.AssignmentId,
		Name:           rec.Mentee.Fullname,
		ProfilePicture: rec.Mentee.ProfilePicture,
		AssignmentURL:  rec.AssignmentURL,
		Grade:          rec.Grade,
		Mentee:         *rec.Mentee.ToMenteeDomain(),
		Assignment:     *rec.Assignment.ToAssignmentDomain(),
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}

func FromMenteeAssignmentDomain(menteeAssignmentDomain *domain.MenteeAssignment) *MenteeAssignment {
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
