package usecase

import (
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type assignmentUseCase struct {
	assignmentRepository domain.AssignmentRepository
	courseRepository     domain.CourseRepository
}

func NewAssignmentUseCase(
	assignmentRepository domain.AssignmentRepository,
	courseRepository domain.CourseRepository,
) domain.AssignmentUseCase {
	return assignmentUseCase{
		assignmentRepository: assignmentRepository,
		courseRepository:     courseRepository,
	}
}

func (au assignmentUseCase) Create(assignmentDomain *domain.Assignment) error {
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

func (au assignmentUseCase) FindById(assignmentId string) (*domain.Assignment, error) {
	assignment, err := au.assignmentRepository.FindById(assignmentId)
	if err != nil {
		return nil, helper.ErrAssignmentNotFound
	}

	return assignment, nil
}

func (au assignmentUseCase) FindByCourseId(courseId string) (*domain.Assignment, error) {
	assignment, err := au.assignmentRepository.FindByCourseId(courseId)
	if err != nil {
		return nil, helper.ErrCourseNotFound
	}

	return assignment, nil
}

func (au assignmentUseCase) Update(assignmentId string, assignmentDomain *domain.Assignment) error {
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

func (au assignmentUseCase) Delete(assignmentId string) error {
	if _, err := au.assignmentRepository.FindById(assignmentId); err != nil {
		return err
	}

	err := au.assignmentRepository.Delete(assignmentId)
	if err != nil {
		return helper.ErrAssignmentNotFound
	}

	return nil
}
