package manage_mentees

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type ManageMenteeController struct {
	manageMenteeUseCase domain.ManageMenteeUseCase
}

func NewManageMenteeController(manageMenteeUseCase domain.ManageMenteeUseCase) *ManageMenteeController {
	return &ManageMenteeController{
		manageMenteeUseCase: manageMenteeUseCase,
	}
}

func (ctrl *ManageMenteeController) HandlerDeleteAccessMentee(c echo.Context) error {
	courseId := c.Param("courseId")
	menteeId := c.Param("menteeId")

	err := ctrl.manageMenteeUseCase.DeleteAccess(menteeId, courseId)

	if err != nil {
		if errors.Is(err, helper.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses hapus akses kursus mentee", nil))
}
