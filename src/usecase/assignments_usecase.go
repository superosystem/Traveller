package usecase

import (
	"github.com/google/uuid"

	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type assignmentUsecase struct {
	assignmentRepository domain.AssignmentRepository
	courseRepository     domain.CourseRepository
}

func NewAssignmentUsecase(assignmentRepository domain.AssignmentRepository, courseRepository domain.CourseRepository) domain.AssignmentUsecase {
	return assignmentUsecase{
		assignmentRepository: assignmentRepository,
		courseRepository:     courseRepository,
	}
}

func (au assignmentUsecase) Create(assignmentDomain *domain.Assignment) error {
	if _, err := au.courseRepository.FindById(assignmentDomain.CourseId); err != nil {
		return err
	}

	id := uuid.NewString()

	assignment := domain.Assignment{
		ID:          id,
		CourseId:    assignmentDomain.CourseId,
		Title:       assignmentDomain.Title,
		Description: assignmentDomain.Description,
	}

	err := au.assignmentRepository.Create(&assignment)

	if err != nil {
		return err
	}

	return nil
}

func (au assignmentUsecase) FindById(assignmentId string) (*domain.Assignment, error) {
	assignment, err := au.assignmentRepository.FindById(assignmentId)

	if err != nil {
		return nil, helper.ErrAssignmentNotFound
	}

	return assignment, nil
}

func (au assignmentUsecase) FindByCourseId(courseId string) (*domain.Assignment, error) {
	assignment, err := au.assignmentRepository.FindByCourseId(courseId)

	if err != nil {
		return nil, helper.ErrCourseNotFound
	}

	return assignment, nil
}

func (au assignmentUsecase) Update(assignmentId string, assignmentDomain *domain.Assignment) error {
	if _, err := au.courseRepository.FindById(assignmentDomain.CourseId); err != nil {
		return err
	}

	_, err := au.assignmentRepository.FindById(assignmentId)

	if err != nil {
		return helper.ErrAssignmentNotFound
	}

	// updatedAssignment := Domain{
	// 	CourseId:    assignmentDomain.CourseId,
	// 	Title:       assignmentDomain.Title,
	// 	Description: assignmentDomain.Description,
	// }

	err = au.assignmentRepository.Update(assignmentId, assignmentDomain)

	if err != nil {
		return helper.ErrAssignmentNotFound
	}

	return nil
}

func (au assignmentUsecase) Delete(assignmentId string) error {
	if _, err := au.assignmentRepository.FindById(assignmentId); err != nil {
		return err
	}

	err := au.assignmentRepository.Delete(assignmentId)

	if err != nil {
		return helper.ErrAssignmentNotFound
	}

	return nil
}
