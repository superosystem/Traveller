package review_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	reviewRepository       mocks.ReviewRepository
	menteeCourseRepository mocks.MenteeCourseRepository
	menteeRepository       mocks.MenteeRepository
	courseRepository       mocks.CourseRepository
	reviewService          domain.ReviewUseCase
	reviewDomain           domain.Review
	menteeCourseDomain     domain.MenteeCourse
	courseDomain           domain.Course
	menteeDomain           domain.Mentee
)

func TestMain(m *testing.M) {
	reviewService = usecase.NewReviewUseCase(&reviewRepository, &menteeCourseRepository, &menteeRepository, &courseRepository)

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

	menteeCourseDomain = domain.MenteeCourse{
		ID:             uuid.NewString(),
		MenteeId:       menteeDomain.ID,
		CourseId:       courseDomain.ID,
		Status:         "completed",
		Reviewed:       true,
		ProgressCount:  10,
		TotalMaterials: 10,
	}

	reviewDomain = domain.Review{
		ID:          uuid.NewString(),
		MenteeId:    menteeDomain.ID,
		CourseId:    courseDomain.ID,
		Description: "test",
		Rating:      int(courseDomain.Rating),
		Reviewed:    menteeCourseDomain.Reviewed,
	}

	m.Run()
}

func TestCreate(t *testing.T) {
	t.Run("Test Create | Success add review", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		reviewRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		menteeCourseRepository.Mock.On("Update", menteeDomain.ID, courseDomain.ID, mock.Anything).Return(nil).Once()

		err := reviewService.Create(&reviewDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Create | Failed add review | Enrollment not found", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(nil, helper.ErrNoEnrolled).Once()

		err := reviewService.Create(&reviewDomain)

		assert.Error(t, err)
	})

	t.Run("Test Create | Failed add review | Error add review", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		reviewRepository.Mock.On("Create", mock.Anything).Return(errors.New("error occurred")).Once()

		err := reviewService.Create(&reviewDomain)

		assert.Error(t, err)
	})

	t.Run("Test Create | Failed add review | Error update mentee course", func(t *testing.T) {
		menteeCourseRepository.Mock.On("CheckEnrollment", menteeDomain.ID, courseDomain.ID).Return(&menteeCourseDomain, nil).Once()

		reviewRepository.Mock.On("Create", mock.Anything).Return(nil)

		menteeCourseRepository.Mock.On("Update", menteeDomain.ID, courseDomain.ID, mock.Anything).Return(errors.New("error occurred")).Once()

		err := reviewService.Create(&reviewDomain)

		assert.Error(t, err)
	})
}

func TestFindByCourse(t *testing.T) {
	t.Run("Test Find By Course | Success find by course", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		reviewRepository.Mock.On("FindByCourse", courseDomain.ID).Return([]domain.Review{reviewDomain}, nil).Once()

		results, err := reviewService.FindByCourse(courseDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find By Course | Failed find by course | Course not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(nil, helper.ErrCourseNotFound).Once()

		results, err := reviewService.FindByCourse(courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test Find By Course | Failed find by course | Review not found", func(t *testing.T) {
		courseRepository.Mock.On("FindById", courseDomain.ID).Return(&courseDomain, nil).Once()

		reviewRepository.Mock.On("FindByCourse", courseDomain.ID).Return(nil, errors.New("Review not found")).Once()

		results, err := reviewService.FindByCourse(courseDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindByMentee(t *testing.T) {
	t.Run("Test Find By Mentee | Success get reviews by mentee", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, courseDomain.Title, menteeCourseDomain.Status).Return(&[]domain.MenteeCourse{menteeCourseDomain}, nil).Once()

		results, err := reviewService.FindByMentee(menteeDomain.ID, courseDomain.Title)

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find By Mentee | Failed get reviews by mentee | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(nil, helper.ErrMenteeNotFound).Once()

		results, err := reviewService.FindByMentee(menteeDomain.ID, courseDomain.Title)

		assert.Error(t, err)
		assert.Empty(t, results)
	})

	t.Run("Test Find By Mentee | Failed get reviews by mentee | Mentee course not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		menteeCourseRepository.Mock.On("FindCoursesByMentee", menteeDomain.ID, courseDomain.Title, menteeCourseDomain.Status).Return(nil, errors.New("not found")).Once()

		results, err := reviewService.FindByMentee(menteeDomain.ID, courseDomain.Title)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}
