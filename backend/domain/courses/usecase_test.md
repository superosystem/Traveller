package courses_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/TrainingSystem/backend/domain/categories"
	_categoryMock "github.com/superosystem/TrainingSystem/backend/domain/categories/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	_courseMock "github.com/superosystem/TrainingSystem/backend/domain/courses/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
	_mentorMock "github.com/superosystem/TrainingSystem/backend/domain/mentors/mocks"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

var (
	courseRepository   _courseMock.Repository
	categoryRepository _categoryMock.CategoryRepositoryMock
	storageClient      helper.StorageConfig
	mentorRepository   _mentorMock.Repository

	courseUsecase courses.Usecase

	courseDomain   courses.Domain
	categoryDomain categories.Domain
	mentorDomain   mentors.Domain
)

func TestMain(m *testing.M) {
	courseUsecase = courses.NewCourseUsecase(&courseRepository, &mentorRepository, &categoryRepository, &storageClient)

	categoryDomain = categories.Domain{
		ID:        uuid.NewString(),
		Name:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mentorDomain = mentors.Domain{
		ID:             uuid.NewString(),
		UserId:         uuid.NewString(),
		Fullname:       "test",
		Email:          "test@gmail.com",
		Phone:          "test",
		Role:           "mentor",
		Jobs:           "test",
		Gender:         "test",
		BirthPlace:     "test",
		BirthDate:      time.Now(),
		Address:        "test",
		ProfilePicture: "test.com",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	courseDomain = courses.Domain{
		ID:          uuid.NewString(),
		MentorId:    mentorDomain.ID,
		CategoryId:  categoryDomain.ID,
		Title:       "test",
		Description: "test",
		Thumbnail:   "test.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Failed create course | Mentor not found", func(t *testing.T) {
		mentorRepository.Mock.On("FindById", mentorDomain.ID).Return(&mentors.Domain{}, helper.ErrMentorNotFound).Once()

		err := courseUsecase.Create(&courseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Create | Failed create course | Error occurred", func(t *testing.T) {
		mentorRepository.Mock.On("FindById", mentorDomain.ID).Return(&mentorDomain, nil).Once()

		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categories.Domain{}, helper.ErrCategoryNotFound).Once()

		err := courseUsecase.Create(&courseDomain)

		assert.Error(t, err)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("Test FindAll | Success get all courses", func(t *testing.T) {
		courseRepository.Mock.On("FindAll", "").Return(&[]courses.Domain{courseDomain}, nil).Once()

		results, err := courseUsecase.FindAll("")

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test FindAll | Failed get all courses | Error occurred", func(t *testing.T) {
		courseRepository.Mock.On("FindAll", "").Return(&[]courses.Domain{}, errors.New("error occurred")).Once()

		results, err := courseUsecase.FindAll("")

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test FindById | Success get course by id", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		result, err := courseUsecase.FindById(courseDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test FindById | Failed get course by id | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courses.Domain{}, helper.ErrCourseNotFound).Once()

		result, err := courseUsecase.FindById(courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestFindByCategory(t *testing.T) {
	t.Run("Test FindByCategory | Success get courses by category", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		courseRepository.Mock.On("FindByCategory", categoryDomain.ID).Return(&[]courses.Domain{courseDomain}, nil).Once()

		results, err := courseUsecase.FindByCategory(categoryDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test FindByCategory | Failed get courses by category | Category not found", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categories.Domain{}, helper.ErrCategoryNotFound).Once()

		results, err := courseUsecase.FindByCategory(categoryDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test FindByCategory | Failed get courses by category | Error occurred", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		courseRepository.Mock.On("FindByCategory", categoryDomain.ID).Return(&[]courses.Domain{}, errors.New("error occurred")).Once()

		results, err := courseUsecase.FindByCategory(categoryDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindByMentor(t *testing.T) {
	t.Run("Test FindByMentor | Success get courses by mentor", func(t *testing.T) {
		mentorRepository.Mock.On("FindById", mentorDomain.ID).Return(&mentorDomain, nil).Once()

		courseRepository.Mock.On("FindByMentor", mentorDomain.ID).Return(&[]courses.Domain{courseDomain}, nil).Once()

		results, err := courseUsecase.FindByMentor(mentorDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test FindByMentor | Failed get courses by mentor | Mentor not found", func(t *testing.T) {
		mentorRepository.Mock.On("FindById", mentorDomain.ID).Return(&mentors.Domain{}, helper.ErrMentorNotFound).Once()

		results, err := courseUsecase.FindByMentor(mentorDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test FindByMentor | Failed get courses by mentor | Error occurred", func(t *testing.T) {
		mentorRepository.Mock.On("FindById", mentorDomain.ID).Return(&mentorDomain, nil).Once()

		courseRepository.Mock.On("FindByMentor", mentorDomain.ID).Return(&[]courses.Domain{}, errors.New("error occurred")).Once()

		results, err := courseUsecase.FindByMentor(mentorDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindByPopular(t *testing.T) {
	t.Run("Test Find By Popular | Success get popular course", func(t *testing.T) {
		courseRepository.Mock.On("FindByPopular").Return([]courses.Domain{courseDomain}, nil).Once()

		results, err := courseUsecase.FindByPopular()

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find By Popular | Failed get popular course | Error occurred", func(t *testing.T) {
		courseRepository.Mock.On("FindByPopular").Return(nil, errors.New("error occurred")).Once()

		results, err := courseUsecase.FindByPopular()

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update course", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		courseRepository.Mock.On("Update", courseDomain.ID, mock.Anything).Return(nil).Once()

		err := courseUsecase.Update(courseDomain.ID, &courseDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update course | Category not found", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categories.Domain{}, helper.ErrCategoryNotFound).Once()

		err := courseUsecase.Update(courseDomain.ID, &courseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update course | Course not found", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courses.Domain{}, helper.ErrCourseNotFound).Once()

		err := courseUsecase.Update(courseDomain.ID, &courseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update course | Error occurred", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		courseRepository.Mock.On("Update", courseDomain.ID, mock.Anything).Return(errors.New("error occurred")).Once()

		err := courseUsecase.Update(courseDomain.ID, &courseDomain)

		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Delete | Success delete course", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		courseRepository.Mock.On("Delete", courseDomain.ID).Return(nil).Once()

		err := courseUsecase.Delete(courseDomain.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed delete course | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courses.Domain{}, helper.ErrCourseNotFound).Once()

		err := courseUsecase.Delete(courseDomain.ID)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed delete course | Error occurred", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		courseRepository.Mock.On("Delete", courseDomain.ID).Return(errors.New("error occurred")).Once()

		err := courseUsecase.Delete(courseDomain.ID)

		assert.Error(t, err)
	})
}
