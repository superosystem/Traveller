package module_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	moduleRepository mocks.ModuleRepository
	courseRepository mocks.CourseRepository
	moduleService    domain.ModuleUsecase
	courseDomain     domain.Course
	moduleDomain     domain.Module
)

func TestMain(m *testing.M) {
	moduleService = usecase.NewModuleUsecase(&moduleRepository, &courseRepository)

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

	moduleDomain = domain.Module{
		ID:        uuid.NewString(),
		CourseId:  courseDomain.ID,
		Title:     "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Success create module", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		moduleRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		err := moduleService.Create(&moduleDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Create | Failed create module | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&domain.Course{}, helper.ErrCourseNotFound).Once()

		err := moduleService.Create(&moduleDomain)

		assert.Error(t, err)
	})

	t.Run("Test Create | Failed create module | Error occurred", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		moduleRepository.Mock.On("Create", mock.Anything).Return(errors.New("error occurred")).Once()

		err := moduleService.Create(&moduleDomain)

		assert.Error(t, err)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test FindById | Success get module by id", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		result, err := moduleService.FindById(moduleDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test FindById | Failed get module by id | Module not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&domain.Module{}, helper.ErrModuleNotFound).Once()

		result, err := moduleService.FindById(moduleDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update module", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		moduleRepository.Mock.On("Update", moduleDomain.ID, &moduleDomain).Return(nil).Once()

		err := moduleService.Update(moduleDomain.ID, &moduleDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update module | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&domain.Course{}, helper.ErrCourseNotFound).Once()

		err := moduleService.Update(moduleDomain.ID, &moduleDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update module | Module not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&domain.Module{}, helper.ErrModuleNotFound).Once()

		err := moduleService.Update(moduleDomain.ID, &moduleDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update module | Error ocurred", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		moduleRepository.Mock.On("Update", moduleDomain.ID, &moduleDomain).Return(errors.New("error occured")).Once()

		err := moduleService.Update(moduleDomain.ID, &moduleDomain)

		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Delete | Success delete module", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		moduleRepository.Mock.On("Delete", moduleDomain.ID).Return(nil).Once()

		err := moduleService.Delete(moduleDomain.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Delete | Failed delete module | Module not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&domain.Module{}, helper.ErrModuleNotFound).Once()

		err := moduleService.Delete(moduleDomain.ID)

		assert.Error(t, err)
	})

	t.Run("Test Delete | Failed delete module | Error not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		moduleRepository.Mock.On("Delete", moduleDomain.ID).Return(errors.New("error occurred")).Once()

		err := moduleService.Delete(moduleDomain.ID)

		assert.Error(t, err)
	})
}
