package response

import (
	"time"

	detailCourse "github.com/superosystem/TrainingSystem/backend/domain/detailCourse"
)

type Course struct {
	CourseId     string     `json:"course_id"`
	CategoryId   string     `json:"category_id"`
	MentorId     string     `json:"mentor_id"`
	Mentor       string     `json:"mentor"`
	Category     string     `json:"category"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Thumbnail    string     `json:"thumbnail"`
	TotalReviews int        `json:"total_reviews"`
	Rating       float32    `json:"rating"`
	Modules      []Module   `json:"modules"`
	Assignment   Assignment `json:"assignment"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Assignment struct {
	AssignmentID string    `json:"assignment_id"`
	CourseId     string    `json:"course_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Module struct {
	ModuleId    string     `json:"module_id"`
	CourseId    string     `json:"course_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Materials   []Material `json:"materials"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Material struct {
	MaterialId  string    `json:"material_id"`
	ModuleId    string    `json:"module_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FullDetailCourse(domain *detailCourse.Domain) *Course {
	modules := make([]Module, len(domain.Modules))

	for i, module := range domain.Modules {
		modules[i].ModuleId = module.ModuleId
		modules[i].CourseId = module.CourseId
		modules[i].Title = module.Title
		modules[i].Description = module.Description
		modules[i].CreatedAt = module.CreatedAt
		modules[i].UpdatedAt = module.UpdatedAt
	}

	for i, module := range modules {
		module.Materials = make([]Material, len(domain.Modules[i].Materials))

		for j, material := range domain.Modules[i].Materials {
			if module.ModuleId == material.ModuleId {
				module.Materials[j].MaterialId = material.MaterialId
				module.Materials[j].ModuleId = material.ModuleId
				module.Materials[j].Title = material.Title
				module.Materials[j].URL = material.URL
				module.Materials[j].Description = material.Description
				module.Materials[j].CreatedAt = material.CreatedAt
				module.Materials[j].UpdatedAt = material.UpdatedAt

				modules[i].Materials = append(modules[i].Materials, module.Materials[j])
			}
		}
	}

	assignment := Assignment{
		AssignmentID: domain.Assignment.ID,
		CourseId:     domain.Assignment.CourseId,
		Title:        domain.Assignment.Title,
		Description:  domain.Assignment.Description,
		CreatedAt:    domain.Assignment.CreatedAt,
		UpdatedAt:    domain.Assignment.UpdatedAt,
	}

	return &Course{
		CourseId:     domain.CourseId,
		CategoryId:   domain.CategoryId,
		MentorId:     domain.MentorId,
		Mentor:       domain.Mentor,
		Category:     domain.Category,
		Title:        domain.Title,
		Description:  domain.Description,
		Thumbnail:    domain.Thumbnail,
		TotalReviews: domain.TotalReviews,
		Rating:       domain.Rating,
		Modules:      modules,
		Assignment:   assignment,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}
}
