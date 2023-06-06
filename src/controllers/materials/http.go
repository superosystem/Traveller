package materials

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/materials/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/materials/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MaterialController struct {
	materialUseCase domain.MaterialUseCase
}

func NewMaterialController(materialUseCase domain.MaterialUseCase) *MaterialController {
	return &MaterialController{
		materialUseCase: materialUseCase,
	}
}

func (ctrl *MaterialController) HandlerCreateMaterial(c echo.Context) error {
	materialInput := request.CreateMaterialInput{}

	materialInput.File, _ = c.FormFile("video")

	if err := c.Bind(&materialInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := materialInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.materialUseCase.Create(materialInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrModuleNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedVideoFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedVideoFile.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan materi", nil))
}

func (ctrl *MaterialController) HandlerFindByIdMaterial(c echo.Context) error {
	materialId := c.Param("materialId")

	material, err := ctrl.materialUseCase.FindById(materialId)

	if err != nil {
		if errors.Is(err, helper.ErrMaterialAssetNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMaterialAssetNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get materi berdasarkan id", response.MaterialDetail(material)))
}

func (ctrl *MaterialController) HandlerUpdateMaterial(c echo.Context) error {
	materialId := c.Param("materialId")
	materialInput := request.UpdateMaterialInput{}

	file, _ := c.FormFile("video")

	if file != nil {
		materialInput.File = file

		if err := c.Bind(&materialInput); err != nil {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
		}
	} else {
		materialInput.ModuleId = c.FormValue("module_id")
		materialInput.Title = c.FormValue("title")
		materialInput.Description = c.FormValue("description")
		materialInput.File = nil
	}

	if err := materialInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.materialUseCase.Update(materialId, materialInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMaterialNotFound.Error()))
		} else if errors.Is(err, helper.ErrMaterialAssetNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMaterialAssetNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedVideoFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedVideoFile.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedVideoFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedVideoFile.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update materi", nil))
}

func (ctrl *MaterialController) HandlerSoftDeleteMaterial(c echo.Context) error {
	materialId := c.Param("materialId")

	err := ctrl.materialUseCase.Delete(materialId)

	if err != nil {
		if errors.Is(err, helper.ErrMaterialAssetNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMaterialAssetNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Materi dihapus", nil))
}

func (ctrl *MaterialController) HandlerSoftDeleteMaterialByModule(c echo.Context) error {
	moduleId := c.Param("moduleId")

	err := ctrl.materialUseCase.Deletes(moduleId)

	if err != nil {
		if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrModuleNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Materi dihapus", nil))
}
