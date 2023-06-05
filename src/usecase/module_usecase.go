package usecase

import (
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type moduleUsecase struct {
	moduleRepository domain.ModuleRepository
	courseRepository domain.CourseRepository
}

func NewModuleUsecase(
	moduleRepository domain.ModuleRepository,
	courseRepository domain.CourseRepository,
) domain.ModuleUsecase {
	return moduleUsecase{
		moduleRepository: moduleRepository,
		courseRepository: courseRepository,
	}
}

func (mu moduleUsecase) Create(moduleDomain *domain.Module) error {
	if _, err := mu.courseRepository.FindById(moduleDomain.CourseId); err != nil {
		return err
	}

	module := domain.Module{
		ID:          uuid.NewString(),
		CourseId:    moduleDomain.CourseId,
		Title:       moduleDomain.Title,
		Description: moduleDomain.Description,
	}

	err := mu.moduleRepository.Create(&module)

	if err != nil {
		return err
	}

	return nil
}

func (mu moduleUsecase) FindById(moduleId string) (*domain.Module, error) {
	module, err := mu.moduleRepository.FindById(moduleId)

	if err != nil {
		return nil, err
	}

	return module, nil
}

func (mu moduleUsecase) Update(moduleId string, moduleDomain *domain.Module) error {
	if _, err := mu.courseRepository.FindById(moduleDomain.CourseId); err != nil {
		return err
	}

	if _, err := mu.moduleRepository.FindById(moduleId); err != nil {
		return err
	}

	err := mu.moduleRepository.Update(moduleId, moduleDomain)

	if err != nil {
		return err
	}

	return nil
}

func (mu moduleUsecase) Delete(moduleId string) error {
	if _, err := mu.moduleRepository.FindById(moduleId); err != nil {
		return err
	}

	err := mu.moduleRepository.Delete(moduleId)

	if err != nil {
		return err
	}

	return nil
}
