package domain

import (
	"mime/multipart"
	"time"

	"github.com/superosystem/trainingsystem-backend/src/helper"
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
	// Create repository create new assignment mentee
	Create(assignmentMenteeDomain *MenteeAssignment) error

	// FindById repository find assignment mentee by id
	FindById(assignmentMenteeId string) (*MenteeAssignment, error)

	// FindByAssignmentID repository find assignment mentee by assignment id
	FindByAssignmentId(assignmentId string, limit int, offset int) ([]MenteeAssignment, int, error)

	// FindByMenteeId repository find assignment mentee by mentee id
	FindByMenteeId(menteeId string) ([]MenteeAssignment, error)

	// FindMenteeAssignmentEnrolled repository find mentee assignment from enrolled course
	FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*MenteeAssignment, error)

	// FindByCourse repository find assignment by course
	FindByCourse(menteeId string, courseId string) (*MenteeAssignment, error)

	// FindByCourses repository find mentee assignments by courses
	FindByCourses(menteeId string, courseIds []string) (*[]MenteeAssignment, error)

	// Update repository update assignment  mentee
	Update(assignmentMenteeId string, assignmentMenteeDomain *MenteeAssignment) error

	// Delete repository delete assignment mentee
	Delete(assignmentMenteeId string) error
}

type MenteeAssignmentUsecase interface {
	// Create usecase create new assignment
	Create(assignmentDomain *MenteeAssignment) error

	// FindById usecase findfind assignment by id
	FindById(assignmentId string) (*MenteeAssignment, error)

	// FindByAssignmentID usecase  find assignment mentee by assignment id
	FindByAssignmentId(assignmentId string, pagination helper.Pagination) (*helper.Pagination, error)

	// FindMenteeAssignmentEnrolled usecase find mentee assignment from enrolled course
	FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*MenteeAssignment, error)

	// FindByMenteeId rusecase find assignment mentee by mentee id
	FindByMenteeId(menteeId string) ([]MenteeAssignment, error)

	// Update usecase update assignment
	Update(assignmentId string, assignmentDomain *MenteeAssignment) error

	// Delete usecase delete assignment
	Delete(assignmentId string) error
}
