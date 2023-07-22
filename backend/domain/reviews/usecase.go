package reviews

import (
	"errors"

	"github.com/google/uuid"
	"github.com/superosystem/TrainingSystem/backend/domain/courses"
	menteeCourses "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses"
	"github.com/superosystem/TrainingSystem/backend/domain/mentees"
)

type reviewUsecase struct {
	reviewRepository       Repository
	menteeCourseRepository menteeCourses.Repository
	menteeRepository       mentees.Repository
	courseRepository       courses.Repository
}

func NewReviewUsecase(
	reviewRepository Repository,
	menteeCourseRepository menteeCourses.Repository,
	menteeRepository mentees.Repository,
	courseRepository courses.Repository,
) Usecase {
	return reviewUsecase{
		reviewRepository:       reviewRepository,
		menteeCourseRepository: menteeCourseRepository,
		menteeRepository:       menteeRepository,
		courseRepository:       courseRepository,
	}
}

func (ru reviewUsecase) Create(reviewDomain *Domain) error {
	if reviewDomain.Rating > 5 || reviewDomain.Rating < 1 {
		return errors.New("Invalid rating value")
	}

	if _, err := ru.menteeCourseRepository.CheckEnrollment(reviewDomain.MenteeId, reviewDomain.CourseId); err != nil {
		return errors.New("Course enrollment not found")
	}

	review := Domain{
		ID:          uuid.NewString(),
		MenteeId:    reviewDomain.MenteeId,
		CourseId:    reviewDomain.CourseId,
		Rating:      reviewDomain.Rating,
		Description: reviewDomain.Description,
	}

	if err := ru.reviewRepository.Create(&review); err != nil {
		return err
	}

	menteeCourse := menteeCourses.Domain{
		Reviewed: true,
	}

	if err := ru.menteeCourseRepository.Update(reviewDomain.MenteeId, reviewDomain.CourseId, &menteeCourse); err != nil {
		return err
	}

	return nil
}

func (ru reviewUsecase) FindByCourse(courseId string) ([]Domain, error) {
	if _, err := ru.courseRepository.FindById(courseId); err != nil {
		return nil, err
	}

	reviews, err := ru.reviewRepository.FindByCourse(courseId)

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (ru reviewUsecase) FindByMentee(menteeId string, title string) ([]Domain, error) {
	if _, err := ru.menteeRepository.FindById(menteeId); err != nil {
		return nil, err
	}

	menteeCourses, err := ru.menteeCourseRepository.FindCoursesByMentee(menteeId, title, "completed")

	if err != nil {
		return nil, err
	}

	reviews := make([]Domain, len(*menteeCourses))

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
