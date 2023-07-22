package mentors

import (
	"errors"
	"net/http"

	"github.com/superosystem/TrainingSystem/backend/controller/mentors/request"
	"github.com/superosystem/TrainingSystem/backend/controller/mentors/response"
	"github.com/superosystem/TrainingSystem/backend/domain/mentors"

	"github.com/labstack/echo/v4"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type MentorController struct {
	mentorUsecase mentors.Usecase
	jwtConfig     *helper.JWTConfig
}

func NewMentorController(mentorUsecase mentors.Usecase, jwtConfig *helper.JWTConfig) *MentorController {
	return &MentorController{
		mentorUsecase: mentorUsecase,
		jwtConfig:     jwtConfig,
	}
}

func (ctrl *MentorController) HandlerRegisterMentor(c echo.Context) error {
	mentorInput := request.MentorRegisterInput{}

	if err := c.Bind(&mentorInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := mentorInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.mentorUsecase.Register(mentorInput.ToDomain())
	if err != nil {
		if errors.Is(err, helper.ErrPasswordLengthInvalid) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordLengthInvalid.Error()))
		} else if errors.Is(err, helper.ErrEmailAlreadyExist) {
			return c.JSON(http.StatusConflict, helper.ConflictResponse(helper.ErrEmailAlreadyExist.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusCreated, helper.SuccessCreatedResponse("Register berhasil", nil))
}

func (ctrl *MentorController) HandlerLoginMentor(c echo.Context) error {
	mentorInput := request.AuthMentorInput{}

	if err := c.Bind(&mentorInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := mentorInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	res, err := ctrl.mentorUsecase.Login(mentorInput.ToDomain())

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

func (ctrl *MentorController) HandlerUpdatePassword(c echo.Context) error {
	mentorId := c.Param("mentorId")
	mentorInput := request.MentorUpdatePassword{}

	if err := c.Bind(&mentorInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}
	mentorInput.UserID = mentorId

	if err := mentorInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.mentorUsecase.UpdatePassword(mentorInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrPasswordLengthInvalid) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordLengthInvalid.Error()))
		} else if errors.Is(err, helper.ErrPasswordNotMatch) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrPasswordNotMatch.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update kata sandi", nil))
}

func (ctrl *MentorController) HandlerFindByID(c echo.Context) error {
	var id string = c.Param("mentorId")

	mentor, err := ctrl.mentorUsecase.FindById(id)

	if err != nil {
		if errors.Is(err, helper.ErrMentorNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMentorNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get mentor berdasarkan id", response.FromDomainUser(mentor)))
}

func (ctrl *MentorController) HandlerProfileMentor(c echo.Context) error {
	token, _ := ctrl.jwtConfig.ExtractToken(c)

	mentor, err := ctrl.mentorUsecase.FindById(token.MentorId)

	if err != nil {
		if errors.Is(err, helper.ErrMentorNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMentorNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get mentor berdasarkan token header", response.FromDomainUser(mentor)))
}

func (ctrl *MentorController) HandlerFindAll(c echo.Context) error {

	mentors, err := ctrl.mentorUsecase.FindAll()

	allMentor := []response.FindMentorAll{}

	for _, mentor := range *mentors {
		allMentor = append(allMentor, *response.FromDomainAll(&mentor))
	}

	if err != nil {
		if errors.Is(err, helper.ErrMentorNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrMentorNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses get semua mentor", allMentor))
}

func (ctrl *MentorController) HandlerUpdateProfile(c echo.Context) error {
	mentorInput := request.MentorUpdateProfile{}

	mentorId := c.Param("mentorId")

	ProfilePictureFile, _ := c.FormFile("profile_picture")

	if ProfilePictureFile != nil {
		mentorInput.ProfilePictureFile = ProfilePictureFile

		if err := c.Bind(&mentorInput); err != nil {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
		}
	} else {
		mentorInput.UserID = c.FormValue("user_id")
		mentorInput.Fullname = c.FormValue("fullname")
		mentorInput.Email = c.FormValue("email")
		mentorInput.Phone = c.FormValue("phone")
		mentorInput.Jobs = c.FormValue("jobs")
		mentorInput.Gender = c.FormValue("gender")
		mentorInput.BirthDate = c.FormValue("birth_date")
		mentorInput.BirthPlace = c.FormValue("birth_place")
		mentorInput.ProfilePictureFile = nil
	}

	if err := mentorInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.mentorUsecase.Update(mentorId, mentorInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrInvalidRequest) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
		} else if errors.Is(err, helper.ErrUnsupportedImageFile) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrUnsupportedImageFile.Error()))
		} else if errors.Is(err, helper.ErrMentorNotFound) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrMentorNotFound.Error()))
		} else if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses update profile", nil))
}

func (ctrl *MentorController) HandlerForgotPassword(c echo.Context) error {
	mentorInput := request.ForgotPasswordInput{}

	if err := c.Bind(&mentorInput); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrInvalidRequest.Error()))
	}

	if err := mentorInput.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
	}

	err := ctrl.mentorUsecase.ForgotPassword(mentorInput.ToDomain())

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, helper.NotFoundResponse(helper.ErrUserNotFound.Error()))
		} else if errors.Is(err, helper.ErrOTPExpired) {
			return c.JSON(http.StatusBadRequest, helper.BadRequestResponse(helper.ErrOTPExpired.Error()))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.InternalServerErrorResponse(helper.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("Sukses ganti kata sandi", nil))
}
