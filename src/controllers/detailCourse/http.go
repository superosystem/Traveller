package detail_course

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/detailCourse/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type DetailCourseController struct {
	detailCourseUseCase domain.DetailCourseUseCase
}

func NewDetailCourseController(detailCourseUseCase domain.DetailCourseUseCase) *DetailCourseController {
	return &DetailCourseController{
		detailCourseUseCase: detailCourseUseCase,
	}
}

func (ctrl *DetailCourseController) HandlerDetailCourse(c echo.Context) error {
	courseId := c.Param("courseId")

	course, err := ctrl.detailCourseUseCase.DetailCourse(courseId)

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get detail kursus", response.FullDetailCourse(course)))
}

func (ctrl *DetailCourseController) HandlerDetailCourseEnrolled(c echo.Context) error {
	menteeId := c.Param("menteeId")
	courseId := c.Param("courseId")

	course, err := ctrl.detailCourseUseCase.DetailCourseEnrolled(menteeId, courseId)

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get detail kursus yang ter-enroll", response.FullDetailCourseEnrolled(course)))
}
