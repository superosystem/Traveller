package category_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	categoryRepository mocks.CategoryRepository
	categoryService    domain.CategoryUseCase
	categoryDomain     domain.Category
)

func TestMain(m *testing.M) {
	categoryService = usecase.NewCategoryUseCase(&categoryRepository)

	categoryDomain = domain.Category{
		ID:        uuid.NewString(),
		Name:      "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Success create category", func(t *testing.T) {
		categoryRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		err := categoryService.Create(&categoryDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Create | Failed create category | Error occurred", func(t *testing.T) {
		categoryRepository.Mock.On("Create", mock.Anything).Return(errors.New("error occurred")).Once()

		err := categoryService.Create(&categoryDomain)

		assert.Error(t, err)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("Test FindAll | Success get all categories", func(t *testing.T) {
		categoryRepository.Mock.On("FindAll").Return(&[]domain.Category{categoryDomain}, nil).Once()

		results, err := categoryService.FindAll()

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test FindAll | Failed get all categories | error occurred", func(t *testing.T) {
		categoryRepository.Mock.On("FindAll").Return(&[]domain.Category{}, errors.New("error occurred")).Once()

		results, err := categoryService.FindAll()

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test FindById | Success get category by id", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		result, err := categoryService.FindById(categoryDomain.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test FindById | Failed get category by id | Category not found", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&domain.Category{}, helper.ErrCategoryNotFound).Once()

		result, err := categoryService.FindById(categoryDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update category", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		categoryRepository.Mock.On("Update", categoryDomain.ID, &categoryDomain).Return(nil).Once()

		err := categoryService.Update(categoryDomain.ID, &categoryDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update category | Category not found", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&domain.Category{}, helper.ErrCategoryNotFound).Once()

		err := categoryService.Update(categoryDomain.ID, &categoryDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update category | Error occurred", func(t *testing.T) {
		categoryRepository.Mock.On("FindById", categoryDomain.ID).Return(&categoryDomain, nil).Once()

		categoryRepository.Mock.On("Update", categoryDomain.ID, &categoryDomain).Return(errors.New("error occurred")).Once()

		err := categoryService.Update(categoryDomain.ID, &categoryDomain)

		assert.Error(t, err)
	})
}
