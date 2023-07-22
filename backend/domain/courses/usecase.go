package courses

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/superosystem/TrainingSystem/backend/domain/categories"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type courseUsecase struct {
	courseRepository   Repository
	mentorRepository   mentors.Repository
	categoryRepository categories.Repository
	storage            *helper.StorageConfig
}

func NewCourseUsecase(
	courseRepository Repository,
	mentorRepository mentors.Repository,
	categoryRepository categories.Repository,
	storage *helper.StorageConfig,
) Usecase {
	return courseUsecase{
		courseRepository:   courseRepository,
		mentorRepository:   mentorRepository,
		categoryRepository: categoryRepository,
		storage:            storage,
	}
}

func (cu courseUsecase) Create(courseDomain *Domain) error {
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

	course := Domain{
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

func (cu courseUsecase) FindAll(keyword string) (*[]Domain, error) {
	courses, err := cu.courseRepository.FindAll(keyword)

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUsecase) FindById(id string) (*Domain, error) {
	course, err := cu.courseRepository.FindById(id)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (cu courseUsecase) FindByCategory(categoryId string) (*[]Domain, error) {
	if _, err := cu.categoryRepository.FindById(categoryId); err != nil {
		return nil, err
	}

	courses, err := cu.courseRepository.FindByCategory(categoryId)

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUsecase) FindByMentor(mentorId string) (*[]Domain, error) {
	if _, err := cu.mentorRepository.FindById(mentorId); err != nil {
		return nil, err
	}

	courses, err := cu.courseRepository.FindByMentor(mentorId)

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUsecase) FindByPopular() ([]Domain, error) {
	courses, err := cu.courseRepository.FindByPopular()

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cu courseUsecase) Update(id string, courseDomain *Domain) error {
	if _, err := cu.categoryRepository.FindById(courseDomain.CategoryId); err != nil {
		return err
	}

	var err error

	var course *Domain
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

	updatedCourse := Domain{
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

func (cu courseUsecase) Delete(id string) error {
	if _, err := cu.courseRepository.FindById(id); err != nil {
		return err
	}

	err := cu.courseRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
