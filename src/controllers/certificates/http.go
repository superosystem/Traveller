package certificates

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type CertificateController struct {
	certificateUseCase domain.CertificateUseCase
}

func NewCertificateController(certificateUseCase domain.CertificateUseCase) *CertificateController {
	return &CertificateController{
		certificateUseCase: certificateUseCase,
	}
}

func (ctrl *CertificateController) HandlerGenerateCert(c echo.Context) error {
	menteeId := c.Param("menteeId")
	courseId := c.Param("courseId")

	data := domain.Certificate{
		MenteeId: menteeId,
		CourseId: courseId,
	}

	cert, err := ctrl.certificateUseCase.GenerateCert(&data)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
	}

	return c.Blob(http.StatusOK, "application/pdf", cert)
}
