package otp

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/controllers/otp/request"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type OTPController struct {
	otpUsecase domain.OtpUsecase
}

func NewOTPController(otpUsecase domain.OtpUsecase) *OTPController {
	return &OTPController{
		otpUsecase: otpUsecase,
	}
}

func (oc OTPController) HandlerSendOTP(c echo.Context) error {
	otpInput := request.OTP{}

	if err := c.Bind(&otpInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := otpInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := oc.otpUsecase.SendOTP(otpInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses kirim OTP ke email", nil))
}

func (oc OTPController) HandlerCheckOTP(c echo.Context) error {
	otpInput := request.CheckOTP{}

	if err := c.Bind(&otpInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := otpInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := oc.otpUsecase.CheckOTP(otpInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else if errors.Is(err, helper.ErrOTPExpired) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrOTPExpired.Error()))
		} else if errors.Is(err, helper.ErrOTPNotMatch) {
			return c.JSON(http.StatusConflict, helper.ConflictResponse(helper.ErrOTPNotMatch.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("OTP matched", nil))
}
