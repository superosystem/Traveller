package certificates

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/TrainingSystem/backend/domain/certificates"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type CertificateController struct {
	certificateUsecase certificates.Usecase
}

func NewCertificateController(certificateUsecase certificates.Usecase) *CertificateController {
	return &CertificateController{
		certificateUsecase: certificateUsecase,
	}
}

func (ctrl *CertificateController) HandlerGenerateCert(c echo.Context) error {
	menteeId := c.Param("menteeId")
	courseId := c.Param("courseId")

	data := certificates.Domain{
		MenteeId: menteeId,
		CourseId: courseId,
	}

	cert, err := ctrl.certificateUsecase.GenerateCert(&data)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.Blob(http.StatusOK, "application/pdf", cert)
}
