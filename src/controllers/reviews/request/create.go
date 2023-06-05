package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type CreateReviewInput struct {
	MenteeId    string `json:"mentee_id" validate:"required"`
	CourseId    string `json:"course_id" validate:"required"`
	Description string `json:"description"`
	Rating      int    `json:"rating" validate:"required"`
}

func (req *CreateReviewInput) ToDomain() *domain.Review {
	return &domain.Review{
		MenteeId:    req.MenteeId,
		CourseId:    req.CourseId,
		Description: req.Description,
		Rating:      req.Rating,
	}
}

func (req *CreateReviewInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		return err
	}

	return nil
}
