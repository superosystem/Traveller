package courses

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/courses/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/courses/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type CourseController struct {
	courseUseCase domain.CourseUseCase
}

func NewCourseController(courseUseCase domain.CourseUseCase) *CourseController {
	return &CourseController{
		courseUseCase: courseUseCase,
	}
}

func (ctrl *CourseController) HandlerCreateCourse(c echo.Context) error {
	courseInput := request.CreateCourseInput{}

	courseInput.Thumbnail, _ = c.FormFile("thumbnail")

	if err := c.Bind(&courseInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := courseInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.courseUseCase.Create(courseInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrMentorNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMentorNotFound.Error()))
		} else if errors.Is(err, helper.ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCategoryNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedImageFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedImageFile.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan kursus baru", nil))
}

func (ctrl *CourseController) HandlerFindAllCourses(c echo.Context) error {
	title := c.QueryParam("keyword")

	coursesDomain, err := ctrl.courseUseCase.FindAll(title)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	var courseResponse []response.FindCourses

	for _, course := range *coursesDomain {
		courseResponse = append(courseResponse, response.AllCourses(&course))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get semua kursus", courseResponse))
}

func (ctrl *CourseController) HandlerFindByIdCourse(c echo.Context) error {
	courseId := c.Param("courseId")

	courseDomain, err := ctrl.courseUseCase.FindById(courseId)

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}

	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get kursus berdasarkan id", response.DetailCourse(courseDomain)))
}

func (ctrl *CourseController) HandlerFindByCategory(c echo.Context) error {
	categoryId := c.Param("categoryId")

	coursesDomain, err := ctrl.courseUseCase.FindByCategory(categoryId)

	if err != nil {
		if errors.Is(err, helper.ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCategoryNotFound.Error()))
		} else if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	var coursesResponse []response.FindCourses

	for _, course := range *coursesDomain {
		coursesResponse = append(coursesResponse, response.AllCourses(&course))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get kursus berdasarkan kategori", coursesResponse))
}

func (ctrl *CourseController) HandlerFindByMentor(c echo.Context) error {
	mentorId := c.Param("mentorId")

	coursesDomain, err := ctrl.courseUseCase.FindByMentor(mentorId)

	if err != nil {
		if errors.Is(err, helper.ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMentorNotFound.Error()))
		} else if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	var coursesResponse []response.FindCourses

	for _, course := range *coursesDomain {
		coursesResponse = append(coursesResponse, response.AllCourses(&course))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get kursus berdasarkan mentor", coursesResponse))
}

func (ctrl *CourseController) HandlerFindByPopular(c echo.Context) error {
	coursesDomain, err := ctrl.courseUseCase.FindByPopular()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	var coursesResponse []response.FindCourses

	for _, course := range coursesDomain {
		coursesResponse = append(coursesResponse, response.AllCourses(&course))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get kursus berdasarkan paling populer", coursesResponse))
}

func (ctrl *CourseController) HandlerUpdateCourse(c echo.Context) error {
	courseId := c.Param("courseId")
	courseInput := request.UpdateCourseInput{}

	// get the thumbnail object file from request body
	thumbnail, _ := c.FormFile("thumbnail")

	if thumbnail != nil {
		courseInput.Thumbnail = thumbnail

		if err := c.Bind(&courseInput); err != nil {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
		}
	} else {
		courseInput.CategoryId = c.FormValue("category_id")
		courseInput.Title = c.FormValue("title")
		courseInput.Description = c.FormValue("description")
		courseInput.Thumbnail = nil
	}

	if err := courseInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.courseUseCase.Update(courseId, courseInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCategoryNotFound.Error()))
		} else if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedImageFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedImageFile.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update kursus", nil))
}

func (ctrl *CourseController) HandlerSoftDeleteCourse(c echo.Context) error {
	courseId := c.Param("courseId")

	err := ctrl.courseUseCase.Delete(courseId)

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Kursus dihapus", nil))
}
