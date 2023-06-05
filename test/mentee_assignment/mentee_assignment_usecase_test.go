package menteeassignment_test

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
	menteeAssignmentRepository mocks.MenteeAssignmentRepository
	assignmentRepository       mocks.AssignmentRepository
	menteeAssignmentService    domain.MenteeAssignmentUsecase
	menteeRepository           mocks.MenteeRepository
	storageClient              config.StorageConfig

	assignmentDomain        domain.Assignment
	menteeDomain            domain.Mentee
	menteeAssignmentDomain  domain.MenteeAssignment
	createMenteeAssignment  domain.MenteeAssignment
	updatedMenteeAssignment domain.MenteeAssignment
)

func TestMain(m *testing.M) {
	menteeAssignmentService = usecase.NewMenteeAssignmentUsecase(&menteeAssignmentRepository, &assignmentRepository, &menteeRepository, &storageClient)

	assignmentDomain = domain.Assignment{
		ID:          uuid.NewString(),
		CourseId:    uuid.NewString(),
		Title:       "test",
		Description: "unit test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	menteeDomain = domain.Mentee{
		ID:       uuid.NewString(),
		UserId:   uuid.NewString(),
		Fullname: "test",
		Phone:    "03536654457",
	}

	menteeAssignmentDomain = domain.MenteeAssignment{
		ID:            uuid.NewString(),
		MenteeId:      menteeDomain.ID,
		AssignmentId:  assignmentDomain.ID,
		AssignmentURL: "test.com",
		Grade:         80,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createMenteeAssignment = domain.MenteeAssignment{
		ID:            uuid.NewString(),
		MenteeId:      menteeDomain.ID,
		AssignmentId:  assignmentDomain.ID,
		AssignmentURL: "test.com",
	}

	updatedMenteeAssignment = domain.MenteeAssignment{
		ID:            uuid.NewString(),
		MenteeId:      menteeDomain.ID,
		AssignmentId:  assignmentDomain.ID,
		AssignmentURL: "test.com",
		PDFfile:       nil,
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Failed create mentee Assignment | Assignmentnot found", func(t *testing.T) {
		assignmentRepository.Mock.On("FindById", assignmentDomain.ID).Return(&domain.Assignment{}, helper.ErrAssignmentNotFound).Once()

		err := menteeAssignmentService.Create(&createMenteeAssignment)

		assert.Error(t, err)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test Find By Id | Success get mentee Assignment by id", func(t *testing.T) {
		menteeAssignmentRepository.Mock.On("FindById", menteeAssignmentDomain.ID).Return(&menteeAssignmentDomain, nil).Once()

		result, err := menteeAssignmentService.FindById(menteeAssignmentDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Find By Id | Failed mentee Assignment not found", func(t *testing.T) {
		menteeAssignmentRepository.Mock.On("FindById", menteeAssignmentDomain.ID).Return(&domain.MenteeAssignment{}, helper.ErrAssignmentMenteeNotFound).Once()

		result, err := menteeAssignmentService.FindById(menteeAssignmentDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update mentee Assignment", func(t *testing.T) {
		assignmentRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		menteeAssignmentRepository.Mock.On("FindById", menteeAssignmentDomain.ID).Return(&menteeAssignmentDomain, nil).Once()

		menteeAssignmentRepository.Mock.On("Update", menteeAssignmentDomain.ID, mock.Anything).Return(nil).Once()

		err := menteeAssignmentService.Update(menteeAssignmentDomain.ID, &updatedMenteeAssignment)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update mentee Assignment | Assignment not found", func(t *testing.T) {
		assignmentRepository.Mock.On("FindById", assignmentDomain.ID).Return(&domain.Assignment{}, helper.ErrAssignmentNotFound).Once()

		err := menteeAssignmentService.Update(menteeAssignmentDomain.ID, &updatedMenteeAssignment)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update mentee Assignment| mentee Assignment not found", func(t *testing.T) {
		assignmentRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		menteeAssignmentRepository.Mock.On("FindById", menteeAssignmentDomain.ID).Return(&domain.MenteeAssignment{}, helper.ErrAssignmentMenteeNotFound).Once()

		err := menteeAssignmentService.Update(menteeAssignmentDomain.ID, &updatedMenteeAssignment)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update mentee Assignment | error occurred", func(t *testing.T) {
		assignmentRepository.Mock.On("FindById", assignmentDomain.ID).Return(&assignmentDomain, nil).Once()

		menteeAssignmentRepository.Mock.On("FindById", menteeAssignmentDomain.ID).Return(&menteeAssignmentDomain, nil).Once()

		menteeAssignmentRepository.Mock.On("Update", menteeAssignmentDomain.ID, mock.Anything).Return(errors.New("error occurred"))

		err := menteeAssignmentService.Update(menteeAssignmentDomain.ID, &updatedMenteeAssignment)

		assert.Error(t, err)
	})
}
