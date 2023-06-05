package menteeprogress_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	menteeProgressRepository mocks.MenteeProgressRepository
	menteeRepository         mocks.MenteeRepository
	courseRepository         mocks.CourseRepository
	materialRepository       mocks.MaterialRepository
	menteeProgressService    domain.MenteeProgressUsecase
	menteeProgressDomain     domain.MenteeProgress
	menteeDomain             domain.Mentee
	courseDomain             domain.Course
	materialDomain           domain.Material
)

func TestMain(m *testing.M) {
	menteeProgressService = usecase.NewMenteeProgressUsecase(
		&menteeProgressRepository,
		&menteeRepository,
		&courseRepository,
		&materialRepository,
	)

	courseDomain = domain.Course{
		ID:          uuid.NewString(),
		MentorId:    uuid.NewString(),
		CategoryId:  uuid.NewString(),
		Title:       "test",
		Description: "test",
		Thumbnail:   "test.com",
	}

	menteeDomain = domain.Mentee{
		ID:             uuid.NewString(),
		UserId:         uuid.NewString(),
		Fullname:       "test",
		Phone:          "test",
		Role:           "mentee",
		Address:        "test",
		ProfilePicture: "test.com",
	}

	materialDomain = domain.Material{
		ID:          uuid.NewString(),
		ModuleId:    uuid.NewString(),
		Title:       "test",
		URL:         "test.com",
		Description: "test",
	}

	menteeProgressDomain = domain.MenteeProgress{
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
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&domain.Mentee{}, helper.ErrMenteeNotFound).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.Error(t, err)
	})

	t.Run("Test Add | Failed add progress | Course not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&domain.Course{}, helper.ErrCourseNotFound).Once()

		err := menteeProgressService.Add(&menteeProgressDomain)

		assert.Error(t, err)
	})

	t.Run("Test Add | Failed add progress | Material not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&domain.Material{}, helper.ErrMaterialNotFound).Once()

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
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&domain.Mentee{}, helper.ErrMenteeNotFound).Once()

		result, err := menteeProgressService.FindMaterialEnrolled(menteeDomain.ID, materialDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test Find Material Enrolled | Failed get material enrolled | Material not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&domain.Material{}, helper.ErrMaterialNotFound).Once()

		result, err := menteeProgressService.FindMaterialEnrolled(menteeDomain.ID, materialDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
