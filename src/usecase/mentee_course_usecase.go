package usecase

import (
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type menteeCourseUseCase struct {
	menteeCourseRepository     domain.MenteeCourseRepository
	menteeRepository           domain.MenteeRepository
	courseRepository           domain.CourseRepository
	materialRepository         domain.MaterialRepository
	menteeProgressRepository   domain.MenteeProgressRepository
	assignmentRepository       domain.AssignmentRepository
	menteeAssignmentRepository domain.MenteeAssignmentRepository
}

func NewMenteeCourseUseCase(
	menteeCourseRepository domain.MenteeCourseRepository,
	menteeRepository domain.MenteeRepository,
	courseRepository domain.CourseRepository,
	materialRepository domain.MaterialRepository,
	menteeProgressRepository domain.MenteeProgressRepository,
	assignmentRepository domain.AssignmentRepository,
	menteeAssignmentRepository domain.MenteeAssignmentRepository,
) domain.MenteeCourseUseCase {
	return menteeCourseUseCase{
		menteeCourseRepository:     menteeCourseRepository,
		menteeRepository:           menteeRepository,
		courseRepository:           courseRepository,
		materialRepository:         materialRepository,
		menteeProgressRepository:   menteeProgressRepository,
		assignmentRepository:       assignmentRepository,
		menteeAssignmentRepository: menteeAssignmentRepository,
	}
}

func (m menteeCourseUseCase) Enroll(menteeCourseDomain *domain.MenteeCourse) error {
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

	menteeCourse := domain.MenteeCourse{
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

func (m menteeCourseUseCase) FindMenteeCourses(menteeId string, title string, status string) (*[]domain.MenteeCourse, error) {
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

func (m menteeCourseUseCase) CheckEnrollment(menteeId string, courseId string) (bool, error) {
	menteeCourseDomain, _ := m.menteeCourseRepository.CheckEnrollment(menteeId, courseId)

	isEnrolled := menteeCourseDomain != nil

	return isEnrolled, nil
}

func (m menteeCourseUseCase) CompleteCourse(menteeId string, courseId string) error {
	if _, err := m.menteeCourseRepository.CheckEnrollment(menteeId, courseId); err != nil {
		return err
	}

	menteeCourse := domain.MenteeCourse{
		Status: "completed",
	}

	err := m.menteeCourseRepository.Update(menteeId, courseId, &menteeCourse)
	if err != nil {
		return err
	}

	return nil
}
