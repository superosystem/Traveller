package domain

import "time"

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
	Enroll(menteeCourseDomain *MenteeCourse) error
	FindCoursesByMentee(menteeId string, title string, status string) (*[]MenteeCourse, error)
	CheckEnrollment(menteeId string, courseId string) (*MenteeCourse, error)
	Update(menteeId string, courseId string, menteeCourseDomain *MenteeCourse) error
	DeleteEnrolledCourse(menteeId string, courseId string) error
}

type MenteeCourseUseCase interface {
	Enroll(menteeCourseDomain *MenteeCourse) error
	FindMenteeCourses(menteeId string, title string, status string) (*[]MenteeCourse, error)
	CheckEnrollment(menteeId string, courseId string) (bool, error)
	CompleteCourse(menteeId string, courseId string) error
}
