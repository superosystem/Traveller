package request

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type UpdateMenteeAssignment struct {
	MenteeID     string                `json:"mentee_id" form:"mentee_id" validate:"required"`
	AssignmentID string                `json:"assignment_id" form:"assignment_id" `
	PDF          *multipart.FileHeader `json:"pdf" form:"pdf" validate:"required"`
}

func (req *UpdateMenteeAssignment) ToDomain() *domain.MenteeAssignment {
	return &domain.MenteeAssignment{
		MenteeId:     req.MenteeID,
		AssignmentId: req.AssignmentID,
		PDFfile:      req.PDF,
	}
}

func (req *UpdateMenteeAssignment) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
