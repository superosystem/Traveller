package response

import (
	"time"

	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
)

type FindMaterialEnrolled struct {
	MenteeId    string    `json:"mentee_id"`
	CourseId    string    `json:"course_id"`
	MaterialId  string    `json:"material_id"`
	Completed   bool      `json:"completed"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DetailMaterialEnrolled(res *menteeProgresses.Domain) *FindMaterialEnrolled {
	return &FindMaterialEnrolled{
		MenteeId:    res.MenteeId,
		CourseId:    res.CourseId,
		MaterialId:  res.MaterialId,
		Completed:   res.Completed,
		Title:       res.Material.Title,
		URL:         res.Material.URL,
		Description: res.Material.Description,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}
