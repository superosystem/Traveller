package domain

import (
	"time"
)

type MenteeCourse struct {
	ID             string
	MenteeId       string
	CourseId       string
	Status         string
	Reviewed       bool
	Mentee         Mentee
	Course         Course
	ProgressCount  int64
	TotalMaterials int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MenteeCourseRepository interface {
	// Enroll repository enroll a course
	Enroll(menteeCourseDomain *MenteeCourse) error

	// FindCoursesByMentee repository find by mentee
	FindCoursesByMentee(menteeId string, title string, status string) (*[]MenteeCourse, error)

	// CheckEnrollment repository check enrollment course mentee
	CheckEnrollment(menteeId string, courseId string) (*MenteeCourse, error)

	// Update repository update course mentee
	Update(menteeId string, courseId string, menteeCourseDomain *MenteeCourse) error

	// DeleteEnrolledCourse delete enrolled course mentee
	DeleteEnrolledCourse(menteeId string, courseId string) error
}

type MenteeCourseUsecase interface {
	// Enroll usecase Enroll usecase mentee enroll course
	Enroll(menteeCourseDomain *MenteeCourse) error

	// FindMenteeCourses usecase find all enrollment courses with title and status
	FindMenteeCourses(menteeId string, title string, status string) (*[]MenteeCourse, error)

	// CheckEnrollment usecase check enrollment course mentee
	CheckEnrollment(menteeId string, courseId string) (bool, error)

	// CompleteCourse usecase to complete the course (update status)
	CompleteCourse(menteeId string, courseId string) error
}
