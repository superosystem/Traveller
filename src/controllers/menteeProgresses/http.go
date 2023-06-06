package mentee_progresses

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/menteeProgresses/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/menteeProgresses/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MenteeProgressController struct {
	menteeProgressUseCase domain.MenteeProgressUseCase
}

func NewMenteeProgressController(menteeProgressUseCase domain.MenteeProgressUseCase) *MenteeProgressController {
	return &MenteeProgressController{
		menteeProgressUseCase: menteeProgressUseCase,
	}
}

func (ctrl *MenteeProgressController) HandlerAddProgress(c echo.Context) error {
	menteeProgressInput := request.AddProgressInput{}

	if err := c.Bind(&menteeProgressInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeProgressInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.menteeProgressUseCase.Add(menteeProgressInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else if errors.Is(err, helper.ErrMaterialAssetNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMaterialAssetNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan progres", nil))
}

func (ctrl *MenteeProgressController) HandlerFindMaterialEnrolled(c echo.Context) error {
	menteeId := c.Param("menteeId")
	materialId := c.Param("materialId")

	progress, err := ctrl.menteeProgressUseCase.FindMaterialEnrolled(menteeId, materialId)

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else if errors.Is(err, helper.ErrMaterialNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get materi", response.DetailMaterialEnrolled(progress)))
}
