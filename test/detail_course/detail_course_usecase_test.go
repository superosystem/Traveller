package detailcourse_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	menteeRepository           mocks.MenteeRepository
	courseRepository           mocks.CourseRepository
	moduleRepository           mocks.ModuleRepository
	materialRepository         mocks.MaterialRepository
	menteeProgressRepository   mocks.MenteeProgressRepository
	assignmentRepository       mocks.AssignmentRepository
	menteeAssignmentRepository mocks.MenteeAssignmentRepository
	menteeCourseRepository     mocks.MenteeCourseRepository
	detailCourseService        domain.DetailCourseUseCase
	menteeDomain               domain.Mentee
	courseDomain               domain.Course
	assignmentDomain           domain.Assignment
	moduleDomain               domain.Module
	materialDomain             domain.Material
	menteeCourseDomain         domain.MenteeCourse
	menteeProgressDomain       domain.MenteeProgress
	menteeAssignmentDomain     domain.MenteeAssignment
)

func TestMain(m *testing.M) {
	detailCourseService = usecase.NewDetailCourseUseCase(
		&menteeRepository,
		&courseRepository,
		&moduleRepository,
		&materialRepository,
		&menteeProgressRepository,
		&assignmentRepository,
		&menteeAssignmentRepository,
		&menteeCourseRepository,
	)

	menteeDomain = domain.Mentee{
		ID:             uuid.NewString(),
		UserId:         uuid.NewString(),
		Fullname:       "test",
		Phone:          "test",
		Role:           "mentee",
		BirthDate:      "test",
		Address:        "test",
		ProfilePicture: "test.com",
	}

	courseDomain = domain.Course{
		ID:           uuid.NewString(),
		MentorId:     uuid.NewString(),
		CategoryId:   uuid.NewString(),
		Title:        "test",
		Description:  "test",
		Thumbnail:    "test.com",
		TotalReviews: 100,
		Rating:       5,
	}

	assignmentDomain = domain.Assignment{
		ID:          uuid.NewString(),
		CourseId:    courseDomain.ID,
		Title:       "test",
		Description: "test",
	}

	moduleDomain = domain.Module{
		ID:          uuid.NewString(),
		CourseId:    courseDomain.ID,
		Title:       "test",
		Description: "test",
	}

	materialDomain = domain.Material{
		ID:          uuid.NewString(),
		CourseId:    courseDomain.ID,
		ModuleId:    moduleDomain.ID,
		Title:       "test",
		URL:         "test.com",
		Description: "test",
	}

	menteeCourseDomain = domain.MenteeCourse{
		ID:             uuid.NewString(),
		MenteeId:       menteeDomain.ID,
		CourseId:       courseDomain.ID,
		Status:         "ongoing",
		Reviewed:       false,
		ProgressCount:  8,
		TotalMaterials: 10,
	}

	menteeProgressDomain = domain.MenteeProgress{
		ID:         uuid.NewString(),
		MenteeId:   menteeDomain.ID,
		CourseId:   courseDomain.ID,
		MaterialId: materialDomain.ID,
		Completed:  false,
	}

	menteeAssignmentDomain = domain.MenteeAssignment{
		ID:             uuid.NewString(),
		MenteeId:       menteeDomain.ID,
		AssignmentId:   assignmentDomain.ID,
		Name:           menteeDomain.Fullname,
		ProfilePicture: menteeDomain.ProfilePicture,
		AssignmentURL:  "test.com",
		Grade:          0,
		Completed:      false,
	}

	m.Run()
}

func TestDetailCourse(t *testing.T) {
	t.Run("Test Detail Course | Success get detail course", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		moduleRepository.Mock.On("FindByCourse", courseDomain.ID).Return([]domain.Module{moduleDomain}, nil).Once()

		assignmentRepository.Mock.On("FindByCourseId", courseDomain.ID).Return(&assignmentDomain, nil).Once()

		materialRepository.Mock.On("FindByModule", []string{moduleDomain.ID}).Return([]domain.Material{materialDomain}, nil).Once()

		result, err := detailCourseService.DetailCourse(courseDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Detail Course | Failed get detail course | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(nil, helper.ErrCourseNotFound).Once()

		result, err := detailCourseService.DetailCourse(courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestDetailCourseEnrolled(t *testing.T) {
	t.Run("Test Detail Course Enrolled | Success get detail course enrolled", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		assignmentRepository.Mock.On("FindByCourseId", courseDomain.ID).Return(&assignmentDomain, nil).Once()

		menteeAssignmentRepository.Mock.On("FindByCourse", menteeDomain.ID, courseDomain.ID).Return(&menteeAssignmentDomain, nil).Once()

		moduleRepository.Mock.On("FindByCourse", courseDomain.ID).Return([]domain.Module{moduleDomain}, nil).Once()

		materialRepository.Mock.On("FindByModule", []string{moduleDomain.ID}).Return([]domain.Material{materialDomain}, nil).Once()

		menteeProgressRepository.Mock.On("FindByMentee", menteeDomain.ID, courseDomain.ID).Return([]domain.MenteeProgress{menteeProgressDomain}, nil).Once()

		materialRepository.Mock.On("CountByCourse", []string{courseDomain.ID}).Return([]int64{menteeCourseDomain.TotalMaterials}, nil).Once()

		menteeProgressRepository.Mock.On("Count", menteeDomain.ID, courseDomain.Title, menteeCourseDomain.Status).Return([]int64{menteeCourseDomain.ProgressCount}, nil).Once()

		result, err := detailCourseService.DetailCourseEnrolled(menteeDomain.ID, courseDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Test Detail Course Enrolled | Failed get detail course enrolled | Course enrollment not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(nil, helper.ErrNoEnrolled).Once()

		result, err := detailCourseService.DetailCourseEnrolled(menteeDomain.ID, courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test Detail Course Enrolled | Failed get detail course enrolled | Course not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(nil, helper.ErrCourseNotFound).Once()

		result, err := detailCourseService.DetailCourseEnrolled(menteeDomain.ID, courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test Detail Course Enrolled | Failed get detail course enrolled | Mentee not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(nil, helper.ErrMenteeNotFound).Once()

		result, err := detailCourseService.DetailCourseEnrolled(menteeDomain.ID, courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
