package domain

import "time"

type DetailCourse struct {
	CourseId       string
	MentorId       string
	CategoryId     string
	Title          string
	Description    string
	Thumbnail      string
	Category       string
	Mentor         string
	TotalReviews   int
	Rating         float32
	Progress       int64
	TotalMaterials int64
	Modules        []DetailModule
	Assignment     DetailAssignment
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type DetailAssignment struct {
	ID          string
	CourseId    string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DetailModule struct {
	ModuleId    string
	CourseId    string
	Title       string
	Description string
	Materials   []DetailMaterial
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DetailMaterial struct {
	MaterialId  string
	ModuleId    string
	Title       string
	URL         string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DetailCourseUseCase interface {
	DetailCourse(courseId string) (*DetailCourse, error)
	DetailCourseEnrolled(menteeId string, courseId string) (*DetailCourse, error)
}
