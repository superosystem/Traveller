package manage_mentees

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type ManageMenteeController struct {
	manageMenteeUsecase domain.ManageMenteeUsecase
}

func NewManageMenteeController(manageMenteeUsecase domain.ManageMenteeUsecase) *ManageMenteeController {
	return &ManageMenteeController{
		manageMenteeUsecase: manageMenteeUsecase,
	}
}

func (ctrl *ManageMenteeController) HandlerDeleteAccessMentee(c echo.Context) error {
	courseId := c.Param("courseId")
	menteeId := c.Param("menteeId")

	err := ctrl.manageMenteeUsecase.DeleteAccess(menteeId, courseId)

	if err != nil {
		if errors.Is(err, helper.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses hapus akses kursus mentee", nil))
}
