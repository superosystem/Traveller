package response

import (
	"time"

	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type FindByIdModule struct {
	ID          string    `json:"id"`
	CourseId    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DetailModule(moduleDomain *domain.Module) *FindByIdModule {
	return &FindByIdModule{
		ID:          moduleDomain.ID,
		CourseId:    moduleDomain.CourseId,
		Title:       moduleDomain.Title,
		Description: moduleDomain.Description,
		CreatedAt:   moduleDomain.CreatedAt,
		UpdatedAt:   moduleDomain.UpdatedAt,
	}
}
