package mentee_assignments

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/TrainingSystem/backend/controller/menteeAssignments/request"
	"github.com/superosystem/TrainingSystem/backend/controller/menteeAssignments/response"
	menteeAssignments "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type AssignmentMenteeController struct {
	assignmentMenteeUsecase menteeAssignments.Usecase
	jwtConfig               *helper.JWTConfig
}

func NewAssignmentsMenteeController(assignmentMenteeUsecase menteeAssignments.Usecase, jwtConfig *helper.JWTConfig) *AssignmentMenteeController {
	return &AssignmentMenteeController{
		assignmentMenteeUsecase: assignmentMenteeUsecase,
		jwtConfig:               jwtConfig,
	}
}

func (ctrl *AssignmentMenteeController) HandlerCreateMenteeAssignment(c echo.Context) error {
	assignmentMenteeInput := request.CreateMenteeAssignment{}

	assignmentMenteeInput.PDF, _ = c.FormFile("pdf")

	if err := c.Bind(&assignmentMenteeInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := assignmentMenteeInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	err := ctrl.assignmentMenteeUsecase.Create(assignmentMenteeInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedAssignmentFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedAssignmentFile.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan tugas mentee", nil))
}

func (ctrl *AssignmentMenteeController) HandlerUpdateMenteeAssignment(c echo.Context) error {
	assignmentMenteeId := c.Param("menteeAssignmentId")
	menteeAssignmentInput := request.UpdateMenteeAssignment{}

	menteeAssignmentInput.PDF, _ = c.FormFile("pdf")

	menteeAssignmentInput.MenteeID = c.FormValue("mentee_id")
	menteeAssignmentInput.AssignmentID = c.FormValue("assignment_id")

	fmt.Println(menteeAssignmentInput)

	if err := c.Bind(&menteeAssignmentInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeAssignmentInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.assignmentMenteeUsecase.Update(assignmentMenteeId, menteeAssignmentInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedAssignmentFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedAssignmentFile.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update tugas mentee", nil))
}

func (ctrl *AssignmentMenteeController) HandlerUpdateGradeMentee(c echo.Context) error {
	id := c.Param("menteeAssignmentId")

	menteeAssignmentInput := request.CreateGrade{}

	if err := c.Bind(&menteeAssignmentInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeAssignmentInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.assignmentMenteeUsecase.Update(id, menteeAssignmentInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentMenteeNotFound.Error()))
		}

		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses Handler nilai", nil))
}

func (ctrl *AssignmentMenteeController) HandlerFindByIdMenteeAssignment(c echo.Context) error {
	id := c.Param("menteeAssignmentId")

	assignmentMentee, err := ctrl.assignmentMenteeUsecase.FindById(id)

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentMenteeNotFound.Error()))
		}

		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get tugas mentee berdasarkan id", response.FromDomain(assignmentMentee)))
}

func (ctrl *AssignmentMenteeController) HandlerFindByAssignmentId(c echo.Context) error {
	id := c.Param("assignmentId")

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	page, _ := strconv.Atoi(c.QueryParam("page"))

	pagination := helper.Pagination{
		Limit: limit,
		Page:  page,
	}

	res, err := ctrl.assignmentMenteeUsecase.FindByAssignmentId(id, pagination)

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentMenteeNotFound.Error()))
		}

		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	var menteeAssignmentResponse []response.AssignmentMentee

	menteeAssignments := res.Result.([]menteeAssignments.Domain)

	for _, menteeAssignment := range menteeAssignments {
		menteeAssignmentResponse = append(menteeAssignmentResponse, response.FromDomain(&menteeAssignment))
	}

	res.Result = menteeAssignmentResponse

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get tugas mentee berdasarkan id tugas", res))
}

func (ctrl *AssignmentMenteeController) HandlerFindByMenteeId(c echo.Context) error {
	token, _ := ctrl.jwtConfig.ExtractToken(c)

	menteeAssignment, err := ctrl.assignmentMenteeUsecase.FindByMenteeId(token.MenteeId)

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}
	var menteeAssignmentResponse []response.AssignmentMentee

	for _, mentee_assignments := range menteeAssignment {
		menteeAssignmentResponse = append(menteeAssignmentResponse, response.FromDomain(&mentee_assignments))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get tugas mentee berdasarkan id mentee", menteeAssignmentResponse))
}

func (ctrl *AssignmentMenteeController) HandlerFindMenteeAssignmentEnrolled(c echo.Context) error {
	menteeId := c.Param("menteeId")
	assignmentId := c.Param("assignmentId")

	menteeAssignment, err := ctrl.assignmentMenteeUsecase.FindMenteeAssignmentEnrolled(menteeId, assignmentId)

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get tugas mentee", response.DetailAssignmentEnrolled(menteeAssignment)))
}

func (ctrl *AssignmentMenteeController) HandlerSoftDeleteMenteeAssignment(c echo.Context) error {
	id := c.Param("menteeAssignmentId")

	err := ctrl.assignmentMenteeUsecase.Delete(id)

	if err != nil {
		if errors.Is(err, helper.ErrAssignmentMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrAssignmentMenteeNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Tugas mentee dihapus", nil))
}
