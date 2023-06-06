package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"path/filepath"
)

type courseUseCase struct {
	courseRepository   domain.CourseRepository
	mentorRepository   domain.MentorRepository
	categoryRepository domain.CategoryRepository
	storage            *config.StorageConfig
}

func NewCourseUseCase(
	courseRepository domain.CourseRepository,
	mentorRepository domain.MentorRepository,
	categoryRepository domain.CategoryRepository,
	storage *config.StorageConfig,
) domain.CourseUseCase {
	return courseUseCase{
		courseRepository:   courseRepository,
		mentorRepository:   mentorRepository,
		categoryRepository: categoryRepository,
		storage:            storage,
	}
}

func (cu courseUseCase) Create(courseDomain *domain.Course) error {
	if _, err := cu.mentorRepository.FindById(courseDomain.MentorId); err != nil {
		return err
	}

	if _, err := cu.categoryRepository.FindById(courseDomain.CategoryId); err != nil {
		return err
	}

	file, err := courseDomain.File.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	extension := filepath.Ext(courseDomain.File.Filename)
	if extension != ".jpg" && extension != ".png" && extension != ".jpeg" {
		return helper.ErrUnsupportedImageFile
	}

	filename, err := helper.GetFilename(courseDomain.File.Filename)
	if err != nil {
		return helper.ErrUnsupportedImageFile
	}

	ctx := context.Background()

	url, err := cu.storage.UploadImage(ctx, filename, file)
	if err != nil {
		return err
	}

	course := domain.Course{
		ID:          uuid.NewString(),
		MentorId:    courseDomain.MentorId,
		CategoryId:  courseDomain.CategoryId,
		Title:       courseDomain.Title,
		Description: courseDomain.Description,
		Thumbnail:   url,
	}

	err = cu.courseRepository.Create(&course)
	if err != nil {
		return err
	}

	return nil
}

func (cu courseUseCase) FindAll(keyword string) (*[]domain.Course, error) {
	courses, err := cu.courseRepository.FindAll(keyword)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUseCase) FindById(id string) (*domain.Course, error) {
	course, err := cu.courseRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (cu courseUseCase) FindByCategory(categoryId string) (*[]domain.Course, error) {
	if _, err := cu.categoryRepository.FindById(categoryId); err != nil {
		return nil, err
	}

	courses, err := cu.courseRepository.FindByCategory(categoryId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUseCase) FindByMentor(mentorId string) (*[]domain.Course, error) {
	if _, err := cu.mentorRepository.FindById(mentorId); err != nil {
		return nil, err
	}

	courses, err := cu.courseRepository.FindByMentor(mentorId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUseCase) FindByPopular() ([]domain.Course, error) {
	courses, err := cu.courseRepository.FindByPopular()
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUseCase) Update(id string, courseDomain *domain.Course) error {
	if _, err := cu.categoryRepository.FindById(courseDomain.CategoryId); err != nil {
		return err
	}

	var err error

	var course *domain.Course
	course, err = cu.courseRepository.FindById(id)
	if err != nil {
		return err
	}

	var url string

	// check if user update the image, do the process
	if courseDomain.File != nil {
		ctx := context.Background()

		if err := cu.storage.DeleteObject(ctx, course.Thumbnail); err != nil {
			return err
		}

		file, err := courseDomain.File.Open()
		if err != nil {
			return err
		}

		defer file.Close()

		extension := filepath.Ext(courseDomain.File.Filename)
		if extension != ".jpg" && extension != ".png" && extension != ".jpeg" {
			return helper.ErrUnsupportedImageFile
		}

		filename, err := helper.GetFilename(courseDomain.File.Filename)
		if err != nil {
			return helper.ErrUnsupportedImageFile
		}

		url, err = cu.storage.UploadImage(ctx, filename, file)
		if err != nil {
			return err
		}
	}

	updatedCourse := domain.Course{
		CategoryId:  courseDomain.CategoryId,
		Title:       courseDomain.Title,
		Description: courseDomain.Description,
		Thumbnail:   url,
	}

	err = cu.courseRepository.Update(id, &updatedCourse)
	if err != nil {
		return err
	}

	return nil
}

func (cu courseUseCase) Delete(id string) error {
	if _, err := cu.courseRepository.FindById(id); err != nil {
		return err
	}

	err := cu.courseRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
