package mentee_courses

import (
	"github.com/google/uuid"
	"github.com/superosystem/TrainingSystem/backend/domain/assignments"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	menteeAssignments "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type menteeCourseUsecase struct {
	menteeCourseRepository     Repository
	menteeRepository           mentees.Repository
	courseRepository           courses.Repository
	materialRepository         materials.Repository
	menteeProgressRepository   menteeProgresses.Repository
	assignmentRepository       assignments.Repository
	menteeAssignmentRepository menteeAssignments.Repository
}

func NewMenteeCourseUsecase(
	menteeCourseRepository Repository,
	menteeRepository mentees.Repository,
	courseRepository courses.Repository,
	materialRepository materials.Repository,
	menteeProgressRepository menteeProgresses.Repository,
	assignmentRepository assignments.Repository,
	menteeAssignmentRepository menteeAssignments.Repository,
) Usecase {
	return menteeCourseUsecase{
		menteeCourseRepository:     menteeCourseRepository,
		menteeRepository:           menteeRepository,
		courseRepository:           courseRepository,
		materialRepository:         materialRepository,
		menteeProgressRepository:   menteeProgressRepository,
		assignmentRepository:       assignmentRepository,
		menteeAssignmentRepository: menteeAssignmentRepository,
	}
}

func (m menteeCourseUsecase) Enroll(menteeCourseDomain *Domain) error {
	if _, err := m.menteeRepository.FindById(menteeCourseDomain.MenteeId); err != nil {
		return err
	}

	if _, err := m.courseRepository.FindById(menteeCourseDomain.CourseId); err != nil {
		return err
	}

	isEnrolled, _ := m.menteeCourseRepository.CheckEnrollment(menteeCourseDomain.MenteeId, menteeCourseDomain.CourseId)

	if isEnrolled != nil {
		return helper.ErrAlreadyEnrolled
	}

	menteeCourseId := uuid.NewString()

	menteeCourse := Domain{
		ID:       menteeCourseId,
		MenteeId: menteeCourseDomain.MenteeId,
		CourseId: menteeCourseDomain.CourseId,
		Status:   "ongoing",
	}

	if err := m.menteeCourseRepository.Enroll(&menteeCourse); err != nil {
		return err
	}

	return nil
}

func (m menteeCourseUsecase) FindMenteeCourses(menteeId string, title string, status string) (*[]Domain, error) {
	menteeCourses, err := m.menteeCourseRepository.FindCoursesByMentee(menteeId, title, status)

	if err != nil {
		return nil, err
	}

	courseIds := []string{}

	for _, course := range *menteeCourses {
		courseIds = append(courseIds, course.CourseId)
	}

	totalMaterials, err := m.materialRepository.CountByCourse(courseIds)

	if err != nil {
		return nil, err
	}

	progresses, err := m.menteeProgressRepository.Count(menteeId, title, status)

	if err != nil {
		return nil, err
	}

	for i, progress := range progresses {
		(*menteeCourses)[i].ProgressCount = progress
	}

	for i, material := range totalMaterials {
		(*menteeCourses)[i].TotalMaterials = material
	}

	assignments, _ := m.assignmentRepository.FindByCourses(courseIds)

	for i := range *menteeCourses {
		for j := range *assignments {
			if (*menteeCourses)[i].CourseId == (*assignments)[j].CourseId {
				(*menteeCourses)[i].TotalMaterials += 1
			}
		}
	}

	menteeAssignments, _ := m.menteeAssignmentRepository.FindByCourses(menteeId, courseIds)

	if menteeAssignments != nil {
		for i := range *menteeCourses {
			for j := range *menteeAssignments {
				if (*menteeCourses)[i].CourseId == (*menteeAssignments)[j].Assignment.CourseId {
					(*menteeCourses)[i].ProgressCount += 1
				}
			}
		}
	}

	return menteeCourses, nil
}

func (m menteeCourseUsecase) CheckEnrollment(menteeId string, courseId string) (bool, error) {
	menteeCourseDomain, _ := m.menteeCourseRepository.CheckEnrollment(menteeId, courseId)

	isEnrolled := menteeCourseDomain != nil

	return isEnrolled, nil
}

func (m menteeCourseUsecase) CompleteCourse(menteeId string, courseId string) error {
	if _, err := m.menteeCourseRepository.CheckEnrollment(menteeId, courseId); err != nil {
		return err
	}

	menteeCourse := Domain{
		Status: "completed",
	}

	err := m.menteeCourseRepository.Update(menteeId, courseId, &menteeCourse)

	if err != nil {
		return err
	}

	return nil
}
