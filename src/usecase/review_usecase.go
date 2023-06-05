package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type reviewUsecase struct {
	reviewRepository       domain.ReviewRepository
	menteeCourseRepository domain.MenteeCourseRepository
	menteeRepository       domain.MenteeRepository
	courseRepository       domain.CourseRepository
}

func NewReviewUsecase(
	reviewRepository domain.ReviewRepository,
	menteeCourseRepository domain.MenteeCourseRepository,
	menteeRepository domain.MenteeRepository,
	courseRepository domain.CourseRepository,
) domain.ReviewUsecase {
	return reviewUsecase{
		reviewRepository:       reviewRepository,
		menteeCourseRepository: menteeCourseRepository,
		menteeRepository:       menteeRepository,
		courseRepository:       courseRepository,
	}
}

func (ru reviewUsecase) Create(reviewDomain *domain.Review) error {
	if reviewDomain.Rating > 5 || reviewDomain.Rating < 1 {
		return errors.New("invalid rating value")
	}

	if _, err := ru.menteeCourseRepository.CheckEnrollment(reviewDomain.MenteeId, reviewDomain.CourseId); err != nil {
		return errors.New("course enrollment not found")
	}

	review := domain.Review{
		ID:          uuid.NewString(),
		MenteeId:    reviewDomain.MenteeId,
		CourseId:    reviewDomain.CourseId,
		Rating:      reviewDomain.Rating,
		Description: reviewDomain.Description,
	}

	if err := ru.reviewRepository.Create(&review); err != nil {
		return err
	}

	menteeCourse := domain.MenteeCourse{
		Reviewed: true,
	}

	if err := ru.menteeCourseRepository.Update(reviewDomain.MenteeId, reviewDomain.CourseId, &menteeCourse); err != nil {
		return err
	}

	return nil
}

func (ru reviewUsecase) FindByCourse(courseId string) ([]domain.Review, error) {
	if _, err := ru.courseRepository.FindById(courseId); err != nil {
		return nil, err
	}

	reviews, err := ru.reviewRepository.FindByCourse(courseId)

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (ru reviewUsecase) FindByMentee(menteeId string, title string) ([]domain.Review, error) {
	if _, err := ru.menteeRepository.FindById(menteeId); err != nil {
		return nil, err
	}

	menteeCourses, err := ru.menteeCourseRepository.FindCoursesByMentee(menteeId, title, "completed")

	if err != nil {
		return nil, err
	}

	reviews := make([]domain.Review, len(*menteeCourses))

	for i, menteeCourse := range *menteeCourses {
		reviews[i].MenteeId = menteeCourse.MenteeId
		reviews[i].CourseId = menteeCourse.CourseId
		reviews[i].Course.Title = menteeCourse.Course.Title
		reviews[i].Course.Mentor.Fullname = menteeCourse.Course.Mentor.Fullname
		reviews[i].Course.Thumbnail = menteeCourse.Course.Thumbnail
		reviews[i].Reviewed = menteeCourse.Reviewed
		reviews[i].CreatedAt = menteeCourse.CreatedAt
		reviews[i].UpdatedAt = menteeCourse.UpdatedAt
	}

	return reviews, nil
}
