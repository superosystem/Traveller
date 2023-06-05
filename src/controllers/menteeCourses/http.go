package mentee_courses

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/controllers/menteeCourses/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/menteeCourses/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type MenteeCourseController struct {
	menteeCourseUsecase domain.MenteeCourseUsecase
}

func NewMenteeCourseController(menteeCourseUsecase domain.MenteeCourseUsecase) *MenteeCourseController {
	return &MenteeCourseController{
		menteeCourseUsecase: menteeCourseUsecase,
	}
}

func (ctrl *MenteeCourseController) HandlerEnrollCourse(c echo.Context) error {
	menteeCourseInput := request.EnrollCourse{}

	if err := c.Bind(&menteeCourseInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeCourseInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.menteeCourseUsecase.Enroll(menteeCourseInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan kursus", nil))
}

func (ctrl *MenteeCourseController) HandlerFindMenteeCourses(c echo.Context) error {
	title := c.QueryParam("keyword")
	status := c.QueryParam("status")
	menteeId := c.Param("menteeId")

	courses, err := ctrl.menteeCourseUsecase.FindMenteeCourses(menteeId, title, status)

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	var menteeCoursesDomain []response.FindMenteeCourses

	for _, course := range *courses {
		menteeCoursesDomain = append(menteeCoursesDomain, *response.MenteeCourses(&course))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get semua kursus mentee", menteeCoursesDomain))
}

func (ctrl *MenteeCourseController) HandlerCheckEnrollmentCourse(c echo.Context) error {
	courseId := c.Param("courseId")
	menteeId := c.Param("menteeId")

	isEnrolled, err := ctrl.menteeCourseUsecase.CheckEnrollment(menteeId, courseId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses ceck status enrollment kursus", map[string]interface{}{
		"status_enrollment": isEnrolled,
	}))
}

func (ctrl *MenteeCourseController) HandlerCompleteCourse(c echo.Context) error {
	courseId := c.Param("courseId")
	menteeId := c.Param("menteeId")

	err := ctrl.menteeCourseUsecase.CompleteCourse(menteeId, courseId)

	if err != nil {
		if errors.Is(err, helper.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses menyelesaikan kursus", nil))
}
