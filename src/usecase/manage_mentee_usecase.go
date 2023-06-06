package usecase

import (
	"context"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type manageMenteeUseCase struct {
	menteeCourse     domain.MenteeCourseRepository
	menteeProgress   domain.MenteeProgressRepository
	menteeAssignment domain.MenteeAssignmentRepository
	storage          *config.StorageConfig
}

func NewManageMenteeUseCase(
	menteeCourse domain.MenteeCourseRepository,
	menteeProgress domain.MenteeProgressRepository,
	menteeAssignment domain.MenteeAssignmentRepository,
	storage *config.StorageConfig,
) domain.ManageMenteeUseCase {
	return manageMenteeUseCase{
		menteeCourse:     menteeCourse,
		menteeProgress:   menteeProgress,
		menteeAssignment: menteeAssignment,
		storage:          storage,
	}
}

func (mm manageMenteeUseCase) DeleteAccess(menteeId string, courseId string) error {
	if _, err := mm.menteeCourse.CheckEnrollment(menteeId, courseId); err != nil {
		return err
	}

	assignment, _ := mm.menteeAssignment.FindByCourse(menteeId, courseId)

	if err := mm.menteeProgress.DeleteMenteeProgressesByCourse(menteeId, courseId); err != nil {
		return err
	}

	if err := mm.menteeCourse.DeleteEnrolledCourse(menteeId, courseId); err != nil {
		return err
	}

	if assignment != nil {
		if err := mm.menteeAssignment.Delete(assignment.ID); err != nil {
			return err
		}

		ctx := context.Background()

		if err := mm.storage.DeleteObject(ctx, assignment.AssignmentURL); err != nil {
			return err
		}
	}

	return nil
}
