package categories

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/controllers/categories/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/categories/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type CategoryController struct {
	categoryUsecase domain.CategoryUsecase
}

func NewCategoryController(categoryUsecase domain.CategoryUsecase) *CategoryController {
	return &CategoryController{
		categoryUsecase: categoryUsecase,
	}
}

func (ctrl *CategoryController) HandlerCreateCategory(c echo.Context) error {
	categoryInput := request.Category{}

	if err := c.Bind(&categoryInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := categoryInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.categoryUsecase.Create(categoryInput.ToDomain())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Sukses menambahkan kategori", nil))
}

func (ctrl *CategoryController) HandlerFindAllCategories(c echo.Context) error {
	categoriesDomain, err := ctrl.categoryUsecase.FindAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	var categoriesResponse []response.Category

	for _, category := range *categoriesDomain {
		categoriesResponse = append(categoriesResponse, response.FromDomain(&category))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get semua kategori", categoriesResponse))
}

func (ctrl *CategoryController) HandlerFindByIdCategory(c echo.Context) error {
	id := c.Param("categoryId")

	category, err := ctrl.categoryUsecase.FindById(id)

	if err != nil {
		if errors.Is(err, helper.ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCategoryNotFound.Error()))
		}

		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get kategori berdasarkan id", response.FromDomain(category)))
}

func (ctrl *CategoryController) HandlerUpdateCategory(c echo.Context) error {
	id := c.Param("categoryId")

	categoryInput := request.Category{}

	if err := c.Bind(&categoryInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := categoryInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.categoryUsecase.Update(id, categoryInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrCategoryNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrCategoryNotFound.Error()))
		}

		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update kategori", nil))
}