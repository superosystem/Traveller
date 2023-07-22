package materials_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	_materialMock "github.com/superosystem/TrainingSystem/backend/domain/materials/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/modules"
	_moduleMock "github.com/superosystem/TrainingSystem/backend/domain/modules/mocks"
	"github.com/superosystem/TrainingSystem/backend/helper"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

var (
	materialRepository _materialMock.Repository
	materialService    materials.Usecase

	moduleRepository _moduleMock.Repository
	storageClient    helper.StorageConfig

	moduleDomain    modules.Domain
	materialDomain  materials.Domain
	createdMaterial materials.Domain
	updatedMaterial materials.Domain
)

func TestMain(m *testing.M) {
	materialService = materials.NewMaterialUsecase(&materialRepository, &moduleRepository, &storageClient)

	moduleDomain = modules.Domain{
		ID:        uuid.NewString(),
		CourseId:  uuid.NewString(),
		Title:     "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	materialDomain = materials.Domain{
		ID:          uuid.NewString(),
		ModuleId:    moduleDomain.ID,
		Title:       "test",
		URL:         "test.com",
		Description: "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdMaterial = materials.Domain{
		ModuleId:    moduleDomain.ID,
		Title:       "test",
		URL:         "test.com",
		Description: "test",
	}

	updatedMaterial = materials.Domain{
		ModuleId:    moduleDomain.ID,
		Title:       "test",
		URL:         "test.com",
		Description: "test",
		File:        nil,
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Failed create material | Module not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&modules.Domain{}, helper.ErrModuleNotFound).Once()

		err := materialService.Create(&createdMaterial)

		assert.Error(t, err)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test Find By Id | Success get material by id", func(t *testing.T) {
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		result, err := materialService.FindById(materialDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Find By Id | Failed material not found", func(t *testing.T) {
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materials.Domain{}, helper.ErrMaterialAssetNotFound).Once()

		result, err := materialService.FindById(materialDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update material", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		materialRepository.Mock.On("Update", materialDomain.ID, mock.Anything).Return(nil).Once()

		err := materialService.Update(materialDomain.ID, &updatedMaterial)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update material | Module not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&modules.Domain{}, helper.ErrModuleNotFound).Once()

		err := materialService.Update(materialDomain.ID, &updatedMaterial)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update material | Material not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materials.Domain{}, helper.ErrMaterialNotFound).Once()

		err := materialService.Update(materialDomain.ID, &updatedMaterial)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update material | error occurred", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		materialRepository.Mock.On("Update", materialDomain.ID, mock.Anything).Return(errors.New("error occurred"))

		err := materialService.Update(materialDomain.ID, &updatedMaterial)

		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Test Delete | Success delete material", func(t *testing.T) {
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		materialRepository.Mock.On("Delete", materialDomain.ID).Return(nil).Once()

		err := materialService.Delete(materialDomain.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Delete | Failed delete material | Material not found", func(t *testing.T) {
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materials.Domain{}, helper.ErrMaterialNotFound).Once()

		err := materialService.Delete(materialDomain.ID)

		assert.Error(t, err)
	})

	t.Run("Test Delete | Failed delete material | Gorm error occurred", func(t *testing.T) {
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&materialDomain, nil).Once()

		materialRepository.Mock.On("Delete", materialDomain.ID).Return(errors.New("error occurred")).Once()

		err := materialService.Delete(materialDomain.ID)

		assert.Error(t, err)
	})
}

func TestDeletes(t *testing.T) {
	t.Run("Test Deletes | Success delete materials", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		materialRepository.Mock.On("Deletes", moduleDomain.ID).Return(nil).Once()

		err := materialService.Deletes(moduleDomain.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Deletes | Failed delete materials | Module not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&modules.Domain{}, helper.ErrModuleNotFound).Once()

		err := materialService.Deletes(moduleDomain.ID)

		assert.Error(t, err)
	})

	t.Run("Test Deletes | Failed delete materials | Error occurred", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		materialRepository.Mock.On("Deletes", moduleDomain.ID).Return(errors.New("error occurred"))

		err := materialService.Deletes(moduleDomain.ID)

		assert.Error(t, err)
	})
}
