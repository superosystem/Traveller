package request

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	menteeAssignments "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
)

type CreateMenteeAssignment struct {
	MenteeID     string                `json:"mentee_id" form:"mentee_id" validate:"required"`
	AssignmentID string                `json:"assignment_id" form:"assignment_id" validate:"required"`
	PDF          *multipart.FileHeader `json:"pdf" form:"pdf" validate:"required"`
}

func (req *CreateMenteeAssignment) ToDomain() *menteeAssignments.Domain {
	return &menteeAssignments.Domain{
		MenteeId:     req.MenteeID,
		AssignmentId: req.AssignmentID,
		PDFfile:      req.PDF,
	}
}

func (req *CreateMenteeAssignment) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
