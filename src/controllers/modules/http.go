package modules

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/modules/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/modules/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type ModuleController struct {
	moduleUseCase domain.ModuleUseCase
}

func NewModuleController(moduleUseCase domain.ModuleUseCase) *ModuleController {
	return &ModuleController{
		moduleUseCase: moduleUseCase,
	}
}

func (ctrl *ModuleController) HandlerCreateModule(c echo.Context) error {
	moduleInput := request.CreateModuleInput{}

	if err := c.Bind(&moduleInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := moduleInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	err := ctrl.moduleUseCase.Create(moduleInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan modul", nil))
}

func (ctrl *ModuleController) HandlerFindByIdModule(c echo.Context) error {
	moduleId := c.Param("moduleId")

	module, err := ctrl.moduleUseCase.FindById(moduleId)

	if err != nil {
		if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrModuleNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get modul berdasarkan id", response.DetailModule(module)))
}

func (ctrl *ModuleController) HandlerUpdateModule(c echo.Context) error {
	moduleId := c.Param("moduleId")
	moduleInput := request.UpdateModuleInput{}

	if err := c.Bind(&moduleInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := moduleInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.moduleUseCase.Update(moduleId, moduleInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrModuleNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update modul", nil))
}

func (ctrl *ModuleController) HandlerDeleteModule(c echo.Context) error {
	moduleId := c.Param("moduleId")

	err := ctrl.moduleUseCase.Delete(moduleId)

	if err != nil {
		if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCourseNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Modul dihapus", nil))
}
