package mentee_progresses

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/TrainingSystem/backend/controller/menteeProgresses/request"
	"github.com/superosystem/TrainingSystem/backend/controller/menteeProgresses/response"
	menteeProgresses "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type MenteeProgressController struct {
	menteeProgressUsecase menteeProgresses.Usecase
}

func NewMenteeProgressController(menteeProgressUsecase menteeProgresses.Usecase) *MenteeProgressController {
	return &MenteeProgressController{
		menteeProgressUsecase: menteeProgressUsecase,
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

	err := ctrl.menteeProgressUsecase.Add(menteeProgressInput.ToDomain())

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

	progress, err := ctrl.menteeProgressUsecase.FindMaterialEnrolled(menteeId, materialId)

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
