package material_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	materialRepository mocks.MaterialRepository
	materialService    domain.MaterialUsecase
	moduleRepository   mocks.ModuleRepository
	storageClient      config.StorageConfig
	moduleDomain       domain.Module
	materialDomain     domain.Material
	createdMaterial    domain.Material
	updatedMaterial    domain.Material
)

func TestMain(m *testing.M) {
	materialService = usecase.NewMaterialUsecase(&materialRepository, &moduleRepository, &storageClient)

	moduleDomain = domain.Module{
		ID:        uuid.NewString(),
		CourseId:  uuid.NewString(),
		Title:     "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	materialDomain = domain.Material{
		ID:          uuid.NewString(),
		ModuleId:    moduleDomain.ID,
		Title:       "test",
		URL:         "test.com",
		Description: "test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdMaterial = domain.Material{
		ModuleId:    moduleDomain.ID,
		Title:       "test",
		URL:         "test.com",
		Description: "test",
	}

	updatedMaterial = domain.Material{
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
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&domain.Module{}, helper.ErrModuleNotFound).Once()

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
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&domain.Material{}, helper.ErrMaterialAssetNotFound).Once()

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
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&domain.Module{}, helper.ErrModuleNotFound).Once()

		err := materialService.Update(materialDomain.ID, &updatedMaterial)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update material | Material not found", func(t *testing.T) {
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&moduleDomain, nil).Once()

		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&domain.Material{}, helper.ErrMaterialNotFound).Once()

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
		materialRepository.Mock.On("FindById", materialDomain.ID).Return(&domain.Material{}, helper.ErrMaterialNotFound).Once()

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
		moduleRepository.Mock.On("FindById", moduleDomain.ID).Return(&domain.Module{}, helper.ErrModuleNotFound).Once()

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
