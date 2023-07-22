package response

import (
	"time"

	"github.com/superosystem/TrainingSystem/backend/domain/modules"
)

type FindByIdModule struct {
	ID          string    `json:"id"`
	CourseId    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DetailModule(moduleDomain *modules.Domain) *FindByIdModule {
	return &FindByIdModule{
		ID:          moduleDomain.ID,
		CourseId:    moduleDomain.CourseId,
		Title:       moduleDomain.Title,
		Description: moduleDomain.Description,
		CreatedAt:   moduleDomain.CreatedAt,
		UpdatedAt:   moduleDomain.UpdatedAt,
	}
}
