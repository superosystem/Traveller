package request

import (
	"github.com/go-playground/validator/v10"
	menteeAssignments "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
)

type CreateGrade struct {
	AssignmentID string `json:"assignment_id" form:"assignment_id" validate:"required"`
	Grade        int    `json:"grade" form:"grade" validate:"required"`
}

func (req *CreateGrade) ToDomain() *menteeAssignments.Domain {
	return &menteeAssignments.Domain{
		AssignmentId: req.AssignmentID,
		Grade:        req.Grade,
	}
}

func (req *CreateGrade) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
