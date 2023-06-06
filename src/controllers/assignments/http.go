package assignments

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/assignments/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/assignments/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type AssignmentController struct {
	assignmentUseCase domain.AssignmentUseCase
}

func NewAssignmentsController(assignmentUseCase domain.AssignmentUseCase) *AssignmentController {
	return &AssignmentController{
		assignmentUseCase: assignmentUseCase,
	}
}

func (ctrl *AssignmentController) HandlerCreateAssignment(c echo.Context) error {
	assignmentInput := request.CreateAssignment{}

	if err := c.Bind(&assignmentInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := assignmentInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.assignmentUseCase.Create(assignmentInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan tugas", nil))
}

func (ctrl *AssignmentController) HandlerFindByIdAssignment(c echo.Context) error {
	assignmentId := c.Param("assignmentId")

	assignment, err := ctrl.assignmentUseCase.FindById(assignmentId)

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get tugas berdasarkan id", response.DetailAssignment(assignment)))
}

func (ctrl *AssignmentController) HandlerFindByCourse(c echo.Context) error {
	courseid := c.Param("courseid")

	assignmentCourse, err := ctrl.assignmentUseCase.FindByCourseId(courseid)

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get tugas berdasarkan id kursus", *response.DetailAssignment(assignmentCourse)))
}

func (ctrl *AssignmentController) HandlerUpdateAssignment(c echo.Context) error {
	assignmentId := c.Param("assignmentId")
	assignmentInput := request.UpdateAssignment{}

	if err := c.Bind(&assignmentInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := assignmentInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.assignmentUseCase.Update(assignmentId, assignmentInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else if errors.Is(err, helper.ErrAssignmentNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update tugas", nil))
}

func (ctrl *AssignmentController) HandlerDeleteAssignment(c echo.Context) error {
	assignmentId := c.Param("assignmentId")

	err := ctrl.assignmentUseCase.Delete(assignmentId)

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Tugas dihapus", nil))
}
