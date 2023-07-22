package response

import (
	"time"

	menteeCourses "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses"
)

type FindMenteeCourses struct {
	ID             string    `json:"id"`
	CourseId       string    `json:"course_id"`
	Mentor         string    `json:"mentor"`
	Title          string    `json:"title"`
	Thumbnail      string    `json:"thumbnail"`
	Progress       int64     `json:"progress"`
	TotalMaterials int64     `json:"total_materials"`
	Status         string    `json:"status"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func MenteeCourses(menteeCourseDomain *menteeCourses.Domain) *FindMenteeCourses {
	return &FindMenteeCourses{
		ID:             menteeCourseDomain.ID,
		CourseId:       menteeCourseDomain.Course.ID,
		Mentor:         menteeCourseDomain.Course.Mentor.Fullname,
		Title:          menteeCourseDomain.Course.Title,
		Thumbnail:      menteeCourseDomain.Course.Thumbnail,
		Progress:       menteeCourseDomain.ProgressCount,
		TotalMaterials: menteeCourseDomain.TotalMaterials,
		Status:         menteeCourseDomain.Status,
		Description:    menteeCourseDomain.Course.Description,
		CreatedAt:      menteeCourseDomain.CreatedAt,
		UpdatedAt:      menteeCourseDomain.UpdatedAt,
	}
}
