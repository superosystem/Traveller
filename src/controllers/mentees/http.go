package mentees

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/controllers/mentees/request"
	"github.com/superosystem/trainingsystem-backend/src/controllers/mentees/response"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type MenteeController struct {
	menteeUseCase domain.MenteeUseCase
	jwtConfig     *config.JWTConfig
}

func NewMenteeController(menteeUseCase domain.MenteeUseCase, jwtConfig *config.JWTConfig) *MenteeController {
	return &MenteeController{
		menteeUseCase: menteeUseCase,
		jwtConfig:     jwtConfig,
	}
}

func (ctrl *MenteeController) HandlerRegisterMentee(c echo.Context) error {
	menteeInput := request.AuthMenteeInput{}

	if err := c.Bind(&menteeInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.menteeUseCase.Register(menteeInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrPasswordLengthInvalid) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordLengthInvalid.Error()))
		} else if errors.Is(err, helper.ErrEmailAlreadyExist) {
			return c.JSON(http.StatusConflict, helper.ConflictResponse(helper.ErrEmailAlreadyExist.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses mengirim OTP ke email", nil))
}

func (ctrl *MenteeController) HandlerVerifyRegisterMentee(c echo.Context) error {
	menteeInput := request.MenteeRegisterInput{}

	if err := c.Bind(&menteeInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.menteeUseCase.VerifyRegister(menteeInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrOTPExpired) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrOTPExpired.Error()))
		} else if errors.Is(err, helper.ErrOTPNotMatch) {
			return c.JSON(http.StatusConflict, helper.ConflictResponse(helper.ErrOTPNotMatch.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Register berhasil", nil))
}

func (ctrl *MenteeController) HandlerLoginMentee(c echo.Context) error {
	menteeInput := request.AuthMenteeInput{}

	if err := c.Bind(&menteeInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	res, err := ctrl.menteeUseCase.Login(menteeInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusConflict, helper.ConflictResponse(helper.ErrAuthenticationFailed.Error()))
		} else if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrPasswordLengthInvalid) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordLengthInvalid.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Login berhasil", res))
}

func (ctrl *MenteeController) HandlerForgotPassword(c echo.Context) error {
	menteeInput := request.ForgotPasswordInput{}

	if err := c.Bind(&menteeInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := menteeInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.menteeUseCase.ForgotPassword(menteeInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrPasswordLengthInvalid) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordLengthInvalid.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else if errors.Is(err, helper.ErrOTPExpired) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrOTPExpired.Error()))
		} else if errors.Is(err, helper.ErrOTPNotMatch) {
			return c.JSON(http.StatusConflict, helper.ConflictResponse(helper.ErrOTPNotMatch.Error()))
		} else if errors.Is(err, helper.ErrPasswordNotMatch) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordNotMatch.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses ganti kata sandi", nil))
}

func (ctrl MenteeController) HandlerFindMenteesByCourse(c echo.Context) error {
	courseId := c.Param("courseId")

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	page, _ := strconv.Atoi(c.QueryParam("page"))

	pagination := helper.Pagination{
		Limit: limit,
		Page:  page,
	}

	res, err := ctrl.menteeUseCase.FindByCourse(courseId, pagination)

	if err != nil {
		if errors.Is(err, helper.ErrCourseNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(err.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	var menteeDomain []response.FindAllMentees

	mentees := res.Result.(*[]domain.Mentee)

	for _, mentee := range *mentees {
		menteeDomain = append(menteeDomain, response.AllMentees(&mentee))
	}

	res.Result = menteeDomain

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get semua mentee berdasarkan kursus", res))
}

func (ctrl *MenteeController) HandlerFindByID(c echo.Context) error {
	var id string = c.Param("menteeId")

	mentee, err := ctrl.menteeUseCase.FindById(id)

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get mentee berdasarkan id", response.FromDomain(mentee)))
}

func (ctrl *MenteeController) HandlerProfileMentee(c echo.Context) error {
	token, _ := ctrl.jwtConfig.ExtractToken(c)

	mentee, err := ctrl.menteeUseCase.FindById(token.MenteeId)

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get mentee berdasarkan token header", response.FromDomain(mentee)))
}

func (ctrl *MenteeController) HandlerFindAll(c echo.Context) error {

	mentees, err := ctrl.menteeUseCase.FindAll()

	allMentees := []response.FindAllMentees{}

	for _, mentee := range *mentees {
		allMentees = append(allMentees, response.AllMentees(&mentee))
	}

	if err != nil {
		if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get semua mentee", allMentees))
}

func (ctrl *MenteeController) HandlerUpdateProfile(c echo.Context) error {
	menteeInput := request.MenteeUpdateProfile{}

	menteeId := c.Param("menteeId")

	ProfilePictureFile, _ := c.FormFile("profile_picture")

	if ProfilePictureFile != nil {
		menteeInput.ProfilePictureFile = ProfilePictureFile

		if err := c.Bind(&menteeInput); err != nil {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
		}
	} else {
		menteeInput.Fullname = c.FormValue("fullname")
		menteeInput.Phone = c.FormValue("phone")
		menteeInput.BirthDate = c.FormValue("birth_date")
		menteeInput.ProfilePictureFile = nil
	}

	if err := menteeInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.menteeUseCase.Update(menteeId, menteeInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrInvalidRequest) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
		} else if errors.Is(err, helper.ErrMenteeNotFound) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrMenteeNotFound.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedImageFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedImageFile.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update profil mentee", nil))
}
