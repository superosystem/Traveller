package domain

import (
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"mime/multipart"
	"time"
)

type MenteeAssignment struct {
	ID             string
	MenteeId       string
	AssignmentId   string
	Name           string
	ProfilePicture string
	AssignmentURL  string
	PDFfile        *multipart.FileHeader
	Grade          int
	Completed      bool
	Mentee         Mentee
	Assignment     Assignment
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MenteeAssignmentRepository interface {
	Create(assignmentMenteeDomain *MenteeAssignment) error
	FindById(assignmentMenteeId string) (*MenteeAssignment, error)
	FindByAssignmentId(assignmentId string, limit int, offset int) ([]MenteeAssignment, int, error)
	FindByMenteeId(menteeId string) ([]MenteeAssignment, error)
	FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*MenteeAssignment, error)
	FindByCourse(menteeId string, courseId string) (*MenteeAssignment, error)
	FindByCourses(menteeId string, courseIds []string) (*[]MenteeAssignment, error)
	Update(assignmentMenteeId string, assignmentMenteeDomain *MenteeAssignment) error
	Delete(assignmentMenteeId string) error
}

type MenteeAssignmentUseCase interface {
	Create(assignmentDomain *MenteeAssignment) error
	FindById(assignmentId string) (*MenteeAssignment, error)
	FindByAssignmentId(assignmentId string, pagination helper.Pagination) (*helper.Pagination, error)
	FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*MenteeAssignment, error)
	FindByMenteeId(menteeId string) ([]MenteeAssignment, error)
	Update(assignmentId string, assignmentDomain *MenteeAssignment) error
	Delete(assignmentId string) error
}
