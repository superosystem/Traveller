package mentee_assignments

import (
	"mime/multipart"
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/assignments"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type Domain struct {
	ID             string
	MenteeId       string
	AssignmentId   string
	Name           string
	ProfilePicture string
	AssignmentURL  string
	PDFfile        *multipart.FileHeader
	Grade          int
	Completed      bool
	Mentee         mentees.Domain
	Assignment     assignments.Domain
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Repository interface {
	// Create repository create new assignment mentee
	Create(assignmentMenteeDomain *Domain) error

	// FindById repository find assignment mentee by id
	FindById(assignmentMenteeId string) (*Domain, error)

	// FindByAssignmentID repository find assignment mentee by assignment id
	FindByAssignmentId(assignmentId string, limit int, offset int) ([]Domain, int, error)

	// FindByMenteeId repository find assignment mentee by mentee id
	FindByMenteeId(menteeId string) ([]Domain, error)

	// FindMenteeAssignmentEnrolled repository find mentee assignment from enrolled course
	FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*Domain, error)

	// FindByCourse repository find assignment by course
	FindByCourse(menteeId string, courseId string) (*Domain, error)

	// FindByCourses repository find mentee assignments by courses
	FindByCourses(menteeId string, courseIds []string) (*[]Domain, error)

	// Update repository update assignment  mentee
	Update(assignmentMenteeId string, assignmentMenteeDomain *Domain) error

	// Delete repository delete assignment mentee
	Delete(assignmentMenteeId string) error
}

type Usecase interface {
	// Create usecase create new assignment
	Create(assignmentDomain *Domain) error

	// FindById usecase findfind assignment by id
	FindById(assignmentId string) (*Domain, error)

	// FindByAssignmentID usecase  find assignment mentee by assignment id
	FindByAssignmentId(assignmentId string, pagination helper.Pagination) (*helper.Pagination, error)

	// FindMenteeAssignmentEnrolled usecase find mentee assignment from enrolled course
	FindMenteeAssignmentEnrolled(menteeId string, assignmentId string) (*Domain, error)

	// FindByMenteeId rusecase find assignment mentee by mentee id
	FindByMenteeId(menteeId string) ([]Domain, error)

	// Update usecase update assignment
	Update(assignmentId string, assignmentDomain *Domain) error

	// Delete usecase delete assignment
	Delete(assignmentId string) error
}
