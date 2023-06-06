package usecase

import (
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type menteeProgressUseCase struct {
	menteeProgressRepository domain.MenteeProgressRepository
	menteeRepository         domain.MenteeRepository
	courseRepository         domain.CourseRepository
	materialRepository       domain.MaterialRepository
}

func NewMenteeProgressUseCase(
	menteeProgressRepository domain.MenteeProgressRepository,
	menteeRepository domain.MenteeRepository,
	courseRepository domain.CourseRepository,
	materialRepository domain.MaterialRepository,
) domain.MenteeProgressUseCase {
	return menteeProgressUseCase{
		menteeProgressRepository: menteeProgressRepository,
		menteeRepository:         menteeRepository,
		courseRepository:         courseRepository,
		materialRepository:       materialRepository,
	}
}

func (m menteeProgressUseCase) Add(menteeProgressDomain *domain.MenteeProgress) error {
	if _, err := m.menteeRepository.FindById(menteeProgressDomain.MenteeId); err != nil {
		return err
	}

	if _, err := m.courseRepository.FindById(menteeProgressDomain.CourseId); err != nil {
		return err
	}

	if _, err := m.materialRepository.FindById(menteeProgressDomain.MaterialId); err != nil {
		return err
	}

	menteeProgress := domain.MenteeProgress{
		ID:         uuid.NewString(),
		MenteeId:   menteeProgressDomain.MenteeId,
		CourseId:   menteeProgressDomain.CourseId,
		MaterialId: menteeProgressDomain.MaterialId,
		Completed:  true,
	}

	err := m.menteeProgressRepository.Add(&menteeProgress)
	if err != nil {
		return err
	}

	return nil
}

func (m menteeProgressUseCase) FindMaterialEnrolled(menteeId string, materialId string) (*domain.MenteeProgress, error) {
	if _, err := m.menteeRepository.FindById(menteeId); err != nil {
		return nil, err
	}

	material, err := m.materialRepository.FindById(materialId)
	if err != nil {
		return nil, err
	}

	progress, _ := m.menteeProgressRepository.FindByMaterial(menteeId, materialId)

	completed := progress != nil

	menteeProgress := domain.MenteeProgress{
		MenteeId:   menteeId,
		CourseId:   material.CourseId,
		MaterialId: materialId,
		Material:   *material,
		Completed:  completed,
		CreatedAt:  material.CreatedAt,
		UpdatedAt:  material.UpdatedAt,
	}

	return &menteeProgress, nil
}
