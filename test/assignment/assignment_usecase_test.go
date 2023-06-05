package assignment_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	mocks "github.com/gusrylmubarok/training-system/ts-backend/src/domain/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	assignmentsRepository mocks.AssignmentRepository
	courseRepository      mocks.CourseRepository
	assignmentService     domain.AssignmentUsecase
	courseDomain          domain.Course
	assignmentDomain      domain.Assignment
)

func TestMain(m *testing.M) {
	assignmentService = usecase.NewAssignmentUsecase(&assignmentsRepository, &courseRepository)

	courseDomain = domain.Course{
		ID:          uuid.NewString(),
		MentorId:    uuid.NewString(),
		CategoryId:  uuid.NewString(),
		Title:       "test",
		Description: "test",
		Thumbnail:   "test.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assignmentDomain = domain.Assignment{
		ID:          uuid.NewString(),
		CourseId:    courseDomain.ID,
		Title:       "test",
		Description: "ini test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Success create assignments", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		assignmentsRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		err := assignmentService.Create(&assignmentDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Create | Failed create assignments | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&domain.Course{}, helper.ErrCourseNotFound).Once()

		err := assignmentService.Create(&assignmentDomain)

		assert.Error(t, err)
	})

	t.Run("Test Create | Failed create assignments | Error occurred", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		assignmentsRepository.Mock.On("Create", mock.Anything).Return(errors.New("error occurred")).Once()

		err := assignmentService.Create(&assignmentDomain)

		assert.Error(t, err)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test FindById | Success get assignments by id", func(t *testing.T) {
		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		result, err := assignmentService.FindById(assignmentDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test FindById | Failed get assignments by id | assignments not found", func(t *testing.T) {
		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&domain.Assignment{}, helper.ErrAssignmentNotFound).Once()

		result, err := assignmentService.FindById(assignmentDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update assignments", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		assignmentsRepository.Mock.On("Update", assignmentDomain.ID, &assignmentDomain).Return(nil).Once()

		err := assignmentService.Update(assignmentDomain.ID, &assignmentDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update assignments | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&domain.Course{}, helper.ErrCourseNotFound).Once()

		err := assignmentService.Update(assignmentDomain.ID, &assignmentDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update assignments| Assignments not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&domain.Assignment{}, helper.ErrAssignmentNotFound).Once()

		err := assignmentService.Update(assignmentDomain.ID, &assignmentDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update assignments | Error ocurred", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		assignmentsRepository.Mock.On("Update", assignmentDomain.ID, &assignmentDomain).Return(errors.New("error occured")).Once()

		err := assignmentService.Update(assignmentDomain.ID, &assignmentDomain)

		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Delete | Success delete assignments", func(t *testing.T) {
		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		assignmentsRepository.Mock.On("Delete", assignmentDomain.ID).Return(nil).Once()

		err := assignmentService.Delete(assignmentDomain.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Delete | Failed delete assignments | assignments not found", func(t *testing.T) {
		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&domain.Assignment{}, helper.ErrAssignmentNotFound).Once()

		err := assignmentService.Delete(assignmentDomain.ID)

		assert.Error(t, err)
	})

	t.Run("Test Delete | Failed delete assignments | Error not found", func(t *testing.T) {
		assignmentsRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		assignmentsRepository.Mock.On("Delete", assignmentDomain.ID).Return(errors.New("error occurred")).Once()

		err := assignmentService.Delete(assignmentDomain.ID)

		assert.Error(t, err)
	})
}
