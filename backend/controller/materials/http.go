package materials

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/TrainingSystem/backend/controller/materials/request"
	"github.com/superosystem/TrainingSystem/backend/controller/materials/response"
	"github.com/superosystem/TrainingSystem/backend/domain/materials"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type MaterialController struct {
	materialUsecase materials.Usecase
}

func NewMaterialController(materialUsecase materials.Usecase) *MaterialController {
	return &MaterialController{
		materialUsecase: materialUsecase,
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

	err := ctrl.materialUsecase.Create(materialInput.ToDomain())

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

	material, err := ctrl.materialUsecase.FindById(materialId)

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

	err := ctrl.materialUsecase.Update(materialId, materialInput.ToDomain())

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

	err := ctrl.materialUsecase.Delete(materialId)

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

	err := ctrl.materialUsecase.Deletes(moduleId)

	if err != nil {
		if errors.Is(err, helper.ErrModuleNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrModuleNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Materi dihapus", nil))
}
