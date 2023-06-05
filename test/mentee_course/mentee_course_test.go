package menteecourse_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	menteeCourseRepository     mocks.MenteeCourseRepository
	courseRepository           mocks.CourseRepository
	menteeRepository           mocks.MenteeRepository
	materialRepository         mocks.MaterialRepository
	menteeProgressRepository   mocks.MenteeProgressRepository
	assignmentRepository       mocks.AssignmentRepository
	menteeAssignmentRepository mocks.MenteeAssignmentRepository

	menteeCourseService domain.MenteeCourseUsecase

	menteeCourseDomain     domain.MenteeCourse
	courseDomain           domain.Course
	menteeDomain           domain.Mentee
	assignmentDomain       domain.Assignment
	menteeAssignmentDomain domain.MenteeAssignment
	progresses             []int64
	totalMaterials         []int64
)

func TestMain(m *testing.M) {
	menteeCourseService = usecase.NewMenteeCourseUsecase(
		&menteeCourseRepository,
		&menteeRepository,
		&courseRepository,
		&materialRepository,
		&menteeProgressRepository,
		&assignmentRepository,
		&menteeAssignmentRepository,
	)

	courseDomain = domain.Course{
		ID:          uuid.NewString(),
		MentorId:    "test",
		CategoryId:  "test",
		Title:       "test",
		Description: "test",
		Thumbnail:   "test.com",
	}

	menteeDomain = domain.Mentee{
		ID:             uuid.NewString(),
		UserId:         "test",
		Fullname:       "test",
		Phone:          "test",
		Role:           "mentee",
		Address:        "test",
		ProfilePicture: "test.com",
	}

	progresses = []int64{5}

	totalMaterials = []int64{10}

	menteeCourseDomain = domain.MenteeCourse{
		ID:             uuid.NewString(),
		MenteeId:       menteeDomain.ID,
		CourseId:       courseDomain.ID,
		Reviewed:       true,
		ProgressCount:  progresses[0],
		TotalMaterials: totalMaterials[0],
		Status:         "ongoing",
	}

	assignmentDomain = domain.Assignment{
		ID:          uuid.NewString(),
		CourseId:    courseDomain.ID,
		Title:       "test",
		Description: "test",
	}

	menteeAssignmentDomain = domain.MenteeAssignment{
		ID:            uuid.NewString(),
		MenteeId:      menteeDomain.ID,
		AssignmentId:  assignmentDomain.ID,
		Name:          menteeDomain.Fullname,
		AssignmentURL: "test.com",
		Grade:         80,
	}

	m.Run()
}

func TestEnroll(t *testing.T) {
	t.Run("Test Enroll | Success enroll course", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(nil, nil).Once()

		menteeCourseRepository.Mock.On("Enroll", mock.Anything).Return(nil).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Enroll | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&domain.Mentee{}, helper.ErrMenteeNotFound).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Enroll | Course not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&domain.Course{}, helper.ErrCourseNotFound).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Enroll | Failed enroll course | Already enrolled course", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&domain.MenteeCourse{}, helper.ErrAlreadyEnrolled).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Enroll | Failed enroll course", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(nil, nil).Once()

		menteeCourseRepository.Mock.On("Enroll", mock.Anything).Return(errors.New("failed enroll course")).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.Error(t, err)
	})
}

func TestFindMenteeCourses(t *testing.T) {
	t.Run("Test Find Mentee Courses | Success get mentee courses", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]domain.MenteeCourse{menteeCourseDomain}, nil).Once()

		materialRepository.Mock.On("CountByCourse", []string{courseDomain.ID}).Return(totalMaterials, nil).Once()

		menteeProgressRepository.Mock.On("Count", menteeDomain.ID, "test", "test").Return(progresses, nil).Once()

		assignmentRepository.Mock.On("FindByCourses", []string{courseDomain.ID}).Return(&[]domain.Assignment{assignmentDomain}, nil).Once()

		menteeAssignmentRepository.Mock.On("FindByCourses", menteeDomain.ID, []string{courseDomain.ID}).Return(&[]domain.MenteeAssignment{menteeAssignmentDomain}, nil).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find Mentee Courses | Failed get mentee courses | Course not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]domain.MenteeCourse{}, helper.ErrCourseNotFound).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test Find Mentee Courses | Failed get mentee courses | error occurred on menteeProgressRepository", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]domain.MenteeCourse{menteeCourseDomain}, nil).Once()

		materialRepository.Mock.On("CountByCourse", []string{courseDomain.ID}).Return(totalMaterials, nil).Once()

		menteeProgressRepository.Mock.On("Count", menteeDomain.ID, "test", "test").Return(nil, errors.New("error occurred")).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test Find Mentee Courses | Failed get mentee courses | error occurred on materialRepository", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]domain.MenteeCourse{menteeCourseDomain}, nil).Once()

		materialRepository.Mock.On("CountByCourse", []string{courseDomain.ID}).Return([]int64{}, errors.New("error occurred")).Once()

		menteeProgressRepository.Mock.On("Count", menteeDomain.ID, "test", "test").Return(progresses, nil).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestCheckEnrollment(t *testing.T) {
	t.Run("Test Check Enrollment | Success check enrollment", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeCourseDomain.MenteeId, menteeCourseDomain.CourseId).Return(&menteeCourseDomain, nil).Once()

		result, err := menteeCourseService.CheckEnrollment(menteeCourseDomain.MenteeId, menteeCourseDomain.CourseId)

		assert.NoError(t, err)
		assert.True(t, result)
	})
}

func TestCompleteCourse(t *testing.T) {
	t.Run("Test Complete Course | Success complete course", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		menteeCourseDomain.Status = "completed"

		menteeCourseRepository.Mock.On("Update", menteeDomain.ID, courseDomain.ID, mock.Anything).Return(nil).Once()

		err := menteeCourseService.CompleteCourse(menteeDomain.ID, courseDomain.ID)

		assert.NoError(t, err)
	})

	t.Run("Test Complete Course | Failed complete course | course enrollment not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(nil, helper.ErrNoEnrolled).Once()

		err := menteeCourseService.CompleteCourse(menteeDomain.ID, courseDomain.ID)

		assert.Error(t, err)
	})

	t.Run("Test Complete Course | Failed complete course | Error update status completion", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		menteeCourseRepository.Mock.On("Update", menteeDomain.ID, courseDomain.ID, mock.Anything).Return(errors.New("error occurred")).Once()

		err := menteeCourseService.CompleteCourse(menteeDomain.ID, courseDomain.ID)

		assert.Error(t, err)
	})
}
