package mentee_progresses

import (
	"github.com/google/uuid"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type menteeProgressUsecase struct {
	menteeProgressRepository Repository
	menteeRepository         mentees.Repository
	courseRepository         courses.Repository
	materialRepository       materials.Repository
}

func NewMenteeProgressUsecase(
	menteeProgressRepository Repository,
	menteeRepository mentees.Repository,
	courseRepository courses.Repository,
	materialRepository materials.Repository,
) Usecase {
	return menteeProgressUsecase{
		menteeProgressRepository: menteeProgressRepository,
		menteeRepository:         menteeRepository,
		courseRepository:         courseRepository,
		materialRepository:       materialRepository,
	}
}

func (m menteeProgressUsecase) Add(menteeProgressDomain *Domain) error {
	if _, err := m.menteeRepository.FindById(menteeProgressDomain.MenteeId); err != nil {
		return err
	}

	if _, err := m.courseRepository.FindById(menteeProgressDomain.CourseId); err != nil {
		return err
	}

	if _, err := m.materialRepository.FindById(menteeProgressDomain.MaterialId); err != nil {
		return err
	}

	menteeProgress := Domain{
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

func (m menteeProgressUsecase) FindMaterialEnrolled(menteeId string, materialId string) (*Domain, error) {
	if _, err := m.menteeRepository.FindById(menteeId); err != nil {
		return nil, err
	}

	material, err := m.materialRepository.FindById(materialId)

	if err != nil {
		return nil, err
	}

	progress, _ := m.menteeProgressRepository.FindByMaterial(menteeId, materialId)

	completed := progress != nil

	menteeProgress := Domain{
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
