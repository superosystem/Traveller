package mentee_courses_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/TrainingSystem/backend/domain/assignments"
	_assignmentMock "github.com/superosystem/TrainingSystem/backend/domain/assignments/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	_courseMock "github.com/superosystem/TrainingSystem/backend/domain/courses/mocks"
	_materialMock "github.com/superosystem/TrainingSystem/backend/domain/materials/mocks"
	menteeAssignments "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
	_menteeAssignmentMock "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments/mocks"
	menteeCourses "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses"
	_menteeCourseMock "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses/mocks"
	_menteeProgressMock "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
	_menteeMock "github.com/superosystem/TrainingSystem/backend/domain/mentees/mocks"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

var (
	menteeCourseRepository     _menteeCourseMock.Repository
	courseRepository           _courseMock.Repository
	menteeRepository           _menteeMock.Repository
	materialRepository         _materialMock.Repository
	menteeProgressRepository   _menteeProgressMock.Repository
	assignmentRepository       _assignmentMock.Repository
	menteeAssignmentRepository _menteeAssignmentMock.Repository

	menteeCourseService menteeCourses.Usecase

	menteeCourseDomain     menteeCourses.Domain
	courseDomain           courses.Domain
	menteeDomain           mentees.Domain
	assignmentDomain       assignments.Domain
	menteeAssignmentDomain menteeAssignments.Domain
	progresses             []int64
	totalMaterials         []int64
)

func TestMain(m *testing.M) {
	menteeCourseService = menteeCourses.NewMenteeCourseUsecase(
		&menteeCourseRepository,
		&menteeRepository,
		&courseRepository,
		&materialRepository,
		&menteeProgressRepository,
		&assignmentRepository,
		&menteeAssignmentRepository,
	)

	courseDomain = courses.Domain{
		ID:          uuid.NewString(),
		MentorId:    "test",
		CategoryId:  "test",
		Title:       "test",
		Description: "test",
		Thumbnail:   "test.com",
	}

	menteeDomain = mentees.Domain{
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

	menteeCourseDomain = menteeCourses.Domain{
		ID:             uuid.NewString(),
		MenteeId:       menteeDomain.ID,
		CourseId:       courseDomain.ID,
		Reviewed:       true,
		ProgressCount:  progresses[0],
		TotalMaterials: totalMaterials[0],
		Status:         "ongoing",
	}

	assignmentDomain = assignments.Domain{
		ID:          uuid.NewString(),
		CourseId:    courseDomain.ID,
		Title:       "test",
		Description: "test",
	}

	menteeAssignmentDomain = menteeAssignments.Domain{
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
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&mentees.Domain{}, helper.ErrMenteeNotFound).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Enroll | Course not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courses.Domain{}, helper.ErrCourseNotFound).Once()

		err := menteeCourseService.Enroll(&menteeCourseDomain)

		assert.Error(t, err)
	})

	t.Run("Test Enroll | Failed enroll course | Already enrolled course", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourses.Domain{}, helper.ErrAlreadyEnrolled).Once()

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
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]menteeCourses.Domain{menteeCourseDomain}, nil).Once()

		materialRepository.Mock.On("CountByCourse", []string{courseDomain.ID}).Return(totalMaterials, nil).Once()

		menteeProgressRepository.Mock.On("Count", menteeDomain.ID, "test", "test").Return(progresses, nil).Once()

		assignmentRepository.Mock.On("FindByCourses", []string{courseDomain.ID}).Return(&[]assignments.Domain{assignmentDomain}, nil).Once()

		menteeAssignmentRepository.Mock.On("FindByCourses", menteeDomain.ID, []string{courseDomain.ID}).Return(&[]menteeAssignments.Domain{menteeAssignmentDomain}, nil).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find Mentee Courses | Failed get mentee courses | Course not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]menteeCourses.Domain{}, helper.ErrCourseNotFound).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test Find Mentee Courses | Failed get mentee courses | error occurred on menteeProgressRepository", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]menteeCourses.Domain{menteeCourseDomain}, nil).Once()

		materialRepository.Mock.On("CountByCourse", []string{courseDomain.ID}).Return(totalMaterials, nil).Once()

		menteeProgressRepository.Mock.On("Count", menteeDomain.ID, "test", "test").Return(nil, errors.New("error occurred")).Once()

		results, err := menteeCourseService.FindMenteeCourses(menteeDomain.ID, "test", "test")

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test Find Mentee Courses | Failed get mentee courses | error occurred on materialRepository", func(t *testing.T) {
		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, "test", "test").Return(&[]menteeCourses.Domain{menteeCourseDomain}, nil).Once()

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
