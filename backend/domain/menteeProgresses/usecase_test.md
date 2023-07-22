package mentee_progresses_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	_courseMock "github.com/superosystem/TrainingSystem/backend/domain/courses/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	_materialMock "github.com/superosystem/TrainingSystem/backend/domain/materials/mocks"
	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
	_menteeProgressMock "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
	_menteeMock "github.com/superosystem/TrainingSystem/backend/domain/mentees/mocks"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

var (
	menteeProgressRepository _menteeProgressMock.Repository
	menteeRepository         _menteeMock.Repository
	courseRepository         _courseMock.Repository
	materialRepository       _materialMock.Repository

	menteeProgressService menteeProgresses.Usecase

	menteeProgressDomain menteeProgresses.Domain
	menteeDomain         mentees.Domain
	courseDomain         courses.Domain
	materialDomain       materials.Domain
)

func TestMain(m *testing.M) {
	menteeProgressService = menteeProgresses.NewMenteeProgressUsecase(
		&menteeProgressRepository,
		&menteeRepository,
		&courseRepository,
		&materialRepository,
	)

	courseDomain = courses.Domain{
		ID:          uuid.NewString(),
		MentorId:    uuid.NewString(),
		CategoryId:  uuid.NewString(),
		Title:       "test",
		Description: "test",
		Thumbnail:   "test.com",
	}

	menteeDomain = mentees.Domain{
		ID:             uuid.NewString(),
		UserId:         uuid.NewString(),
		Fullname:       "test",
		Phone:          "test",
		Role:           "mentee",
		Address:        "test",
		ProfilePicture: "test.com",
	}

	materialDomain = materials.Domain{
		ID:          uuid.NewString(),
		ModuleId:    uuid.NewString(),
		Title:       "test",
		URL:         "test.com",
		Description: "test",
	}

	menteeProgressDomain = menteeProgresses.Domain{
		ID:         uuid.NewString(),
		MenteeId:   menteeDomain.ID,
		CourseId:   courseDomain.ID,
		MaterialId: materialDomain.ID,
	}

	m.Run()
}

func TestAdd(t *testing.T) {
	t.Run("Test Add | Success add progress", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		menteeProgressRepository.Mock.On("Add", mock.Anything).Return(nil).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Add | Failed add progress | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&mentees.Domain{}, helper.ErrMenteeNotFound).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.Error(t, err)
	})

	t.Run("Test Add | Failed add progress | Course not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courses.Domain{}, helper.ErrCourseNotFound).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.Error(t, err)
	})

	t.Run("Test Add | Failed add progress | Material not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materials.Domain{}, helper.ErrMaterialNotFound).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.Error(t, err)
	})

	t.Run("Test Add | Failed add progress | Error occurred", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		menteeProgressRepository.Mock.On("Add", mock.Anything).Return(errors.New("failed add progress")).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.Error(t, err)
	})
}

func TestFindMaterialEnrolled(t *testing.T) {
	t.Run("Test Find Material Enrolled | Success get material enrolled", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		menteeProgressRepository.Mock.On("FindByMaterial", menteeDomain.ID, materialDomain.ID).Return(&menteeProgressDomain, nil).Once()

		result, err := menteeProgressService.FindMaterialEnrolled(menteeDomain.ID, materialDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Find Material Enrolled | Failed get material enrolled | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&mentees.Domain{}, helper.ErrMenteeNotFound).Once()

		result, err := menteeProgressService.FindMaterialEnrolled(menteeDomain.ID, materialDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test Find Material Enrolled | Failed get material enrolled | Material not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materials.Domain{}, helper.ErrMaterialNotFound).Once()

		result, err := menteeProgressService.FindMaterialEnrolled(menteeDomain.ID, materialDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
