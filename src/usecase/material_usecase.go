package usecase

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type materialUsecase struct {
	materialRepository domain.MaterialRepository
	moduleRepository   domain.ModuleRepository
	storage            *config.StorageConfig
}

func NewMaterialUsecase(
	materialRepository domain.MaterialRepository,
	moduleRepository domain.ModuleRepository,
	storage *config.StorageConfig,
) domain.MaterialUsecase {
	return materialUsecase{
		materialRepository: materialRepository,
		moduleRepository:   moduleRepository,
		storage:            storage,
	}
}

func (mu materialUsecase) Create(materialDomain *domain.Material) error {
	if _, err := mu.moduleRepository.FindById(materialDomain.ModuleId); err != nil {
		return err
	}

	file, err := materialDomain.File.Open()

	if err != nil {
		return err
	}

	defer file.Close()

	extension := filepath.Ext(materialDomain.File.Filename)

	if extension != ".mp4" && extension != ".mkv" {
		return helper.ErrUnsupportedVideoFile
	}

	filename, err := helper.GetFilename(materialDomain.File.Filename)

	if err != nil {
		return helper.ErrUnsupportedVideoFile
	}

	ctx := context.Background()

	url, err := mu.storage.UploadVideo(ctx, filename, file)

	if err != nil {
		return err
	}

	material := domain.Material{
		ID:          uuid.NewString(),
		ModuleId:    materialDomain.ModuleId,
		Title:       materialDomain.Title,
		URL:         url,
		Description: materialDomain.Description,
	}

	if err := mu.materialRepository.Create(&material); err != nil {
		return err
	}

	return nil
}

func (mu materialUsecase) FindById(materialId string) (*domain.Material, error) {
	material, err := mu.materialRepository.FindById(materialId)

	if err != nil {
		return nil, err
	}

	return material, nil
}

func (mu materialUsecase) Update(materialId string, materialDomain *domain.Material) error {
	if _, err := mu.moduleRepository.FindById(materialDomain.ModuleId); err != nil {
		return err
	}

	material, err := mu.materialRepository.FindById(materialId)

	if err != nil {
		return err
	}

	var url string

	if materialDomain.File != nil {
		ctx := context.Background()

		err := mu.storage.DeleteObject(ctx, material.URL)

		if err != nil {
			return err
		}

		file, err := materialDomain.File.Open()

		if err != nil {
			return err
		}

		defer file.Close()

		extension := filepath.Ext(materialDomain.File.Filename)

		if extension != ".mp4" && extension != ".mkv" {
			return helper.ErrUnsupportedVideoFile
		}

		filename, err := helper.GetFilename(materialDomain.File.Filename)

		if err != nil {
			return helper.ErrUnsupportedVideoFile
		}

		url, err = mu.storage.UploadVideo(ctx, filename, file)

		if err != nil {
			return err
		}
	}

	updatedMaterial := domain.Material{
		ModuleId:    materialDomain.ModuleId,
		Title:       materialDomain.Title,
		URL:         url,
		Description: materialDomain.Description,
	}

	if err := mu.materialRepository.Update(materialId, &updatedMaterial); err != nil {
		return err
	}

	return nil
}

func (mu materialUsecase) Delete(materialId string) error {
	if _, err := mu.materialRepository.FindById(materialId); err != nil {
		return err
	}

	err := mu.materialRepository.Delete(materialId)

	if err != nil {
		return err
	}

	return nil
}

func (mu materialUsecase) Deletes(moduleId string) error {
	if _, err := mu.moduleRepository.FindById(moduleId); err != nil {
		return err
	}

	err := mu.materialRepository.Deletes(moduleId)

	if err != nil {
		return err
	}

	return nil
}
