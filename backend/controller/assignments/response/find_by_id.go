package response

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/assignments"
)

type FindByIdAssignments struct {
	ID          string    `json:"id"`
	CourseId    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DetailAssignment(assignmentDomain *assignments.Domain) *FindByIdAssignments {
	return &FindByIdAssignments{
		ID:          assignmentDomain.ID,
		CourseId:    assignmentDomain.CourseId,
		Title:       assignmentDomain.Title,
		Description: assignmentDomain.Description,
		CreatedAt:   assignmentDomain.CreatedAt,
		UpdatedAt:   assignmentDomain.UpdatedAt,
	}
}
