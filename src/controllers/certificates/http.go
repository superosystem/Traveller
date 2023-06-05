package certificates

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type CertificateController struct {
	certificateUsecase domain.CertificateUsecase
}

func NewCertificateController(certificateUsecase domain.CertificateUsecase) *CertificateController {
	return &CertificateController{
		certificateUsecase: certificateUsecase,
	}
}

func (ctrl *CertificateController) HandlerGenerateCert(c echo.Context) error {
	menteeId := c.Param("menteeId")
	courseId := c.Param("courseId")

	data := domain.Certificate{
		MenteeId: menteeId,
		CourseId: courseId,
	}

	cert, err := ctrl.certificateUsecase.GenerateCert(&data)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.Blob(http.StatusOK, "application/pdf", cert)
}
