package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"math"
	"path/filepath"
	"time"
)

type menteeUseCase struct {
	menteeRepository domain.MenteeRepository
	userRepository   domain.UserRepository
	otpRepository    domain.OtpRepository
	jwtConfig        *config.JWTConfig
	mailerConfig     *config.MailerConfig
	storage          *config.StorageConfig
}

func NewMenteeUseCase(
	menteeRepository domain.MenteeRepository,
	userRepository domain.UserRepository,
	otpRepository domain.OtpRepository,
	jwtConfig *config.JWTConfig,
	mailerConfig *config.MailerConfig,
	storage *config.StorageConfig,
) domain.MenteeUseCase {
	return menteeUseCase{
		menteeRepository: menteeRepository,
		userRepository:   userRepository,
		otpRepository:    otpRepository,
		jwtConfig:        jwtConfig,
		mailerConfig:     mailerConfig,
		storage:          storage,
	}
}

func (m menteeUseCase) Register(menteeAuth *domain.MenteeAuth) error {
	if len(menteeAuth.Password) < 6 {
		return helper.ErrPasswordLengthInvalid
	}

	user, _ := m.userRepository.FindByEmail(menteeAuth.Email)
	if user != nil {
		return helper.ErrEmailAlreadyExist
	}

	newOTP := helper.GenerateOTP(4)

	var err error

	ctx := context.Background()

	err = m.otpRepository.Save(ctx, menteeAuth.Email, newOTP, helper.TIME_TO_LIVE)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("OTP: %s", newOTP)
	subject := "Verification Registering Training System"

	_ = m.mailerConfig.SendMail(menteeAuth.Email, subject, message)

	return nil
}

func (m menteeUseCase) VerifyRegister(menteeDomain *domain.MenteeRegister) error {
	if len(menteeDomain.Password) < 6 {
		return helper.ErrPasswordLengthInvalid
	}

	userDomain, _ := m.userRepository.FindByEmail(menteeDomain.Email)

	if userDomain != nil {
		return helper.ErrEmailAlreadyExist
	}

	var err error

	ctx := context.Background()
	var validOTP string

	validOTP, err = m.otpRepository.Get(ctx, menteeDomain.Email)
	if err != nil {
		return err
	}

	if validOTP != menteeDomain.OTP {
		return helper.ErrOTPNotMatch
	}

	userId := uuid.NewString()
	hashedPassword := helper.HashPassword(menteeDomain.Password)

	user := domain.User{
		ID:        userId,
		Email:     menteeDomain.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = m.userRepository.Create(&user)
	if err != nil {
		return err
	}

	menteeId := uuid.NewString()

	mentee := domain.Mentee{
		ID:        menteeId,
		UserId:    userId,
		Fullname:  menteeDomain.Fullname,
		Phone:     menteeDomain.Phone,
		Role:      "mentee",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = m.menteeRepository.Create(&mentee)
	if err != nil {
		return err
	}

	return nil
}

func (m menteeUseCase) ForgotPassword(forgotPassword *domain.MenteeForgotPassword) error {
	var err error

	if len(forgotPassword.Password) < 6 || len(forgotPassword.RepeatedPassword) < 6 {
		return helper.ErrPasswordLengthInvalid
	}

	var user *domain.User
	user, err = m.userRepository.FindByEmail(forgotPassword.Email)
	if err != nil {
		return err
	}

	ctx := context.Background()
	var result string

	result, err = m.otpRepository.Get(ctx, forgotPassword.Email)
	if err != nil {
		return err
	}

	if result != forgotPassword.OTP {
		return helper.ErrOTPNotMatch
	}

	if forgotPassword.Password != forgotPassword.RepeatedPassword {
		return helper.ErrPasswordNotMatch
	}

	hashPassword := helper.HashPassword(forgotPassword.RepeatedPassword)

	updatedUser := domain.User{
		Password: hashPassword,
	}

	err = m.userRepository.Update(user.ID, &updatedUser)
	if err != nil {
		return err
	}

	return nil
}

func (m menteeUseCase) Login(menteeAuth *domain.MenteeAuth) (interface{}, error) {
	if len(menteeAuth.Password) < 6 {
		return nil, helper.ErrPasswordLengthInvalid
	}

	var err error

	var user *domain.User

	user, err = m.userRepository.FindByEmail(menteeAuth.Email)
	if err != nil {
		return nil, err
	}

	ok := helper.ComparePassword(user.Password, menteeAuth.Password)
	if !ok {
		return nil, helper.ErrUserNotFound
	}

	var mentee *domain.Mentee
	mentee, err = m.menteeRepository.FindByIdUser(user.ID)
	if err != nil {
		return nil, err
	}

	var token string
	exp := time.Now().Add(6 * time.Hour)

	token, err = m.jwtConfig.GenerateToken(user.ID, mentee.ID, mentee.Role, exp)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"token":   token,
		"expires": exp,
	}

	return data, nil
}

func (m menteeUseCase) FindAll() (*[]domain.Mentee, error) {
	var err error

	mentees, err := m.menteeRepository.FindAll()
	if err != nil {
		if err == helper.ErrMenteeNotFound {
			return nil, helper.ErrMenteeNotFound
		}

		return nil, helper.ErrInternalServerError
	}

	return mentees, nil
}

func (m menteeUseCase) FindById(id string) (*domain.Mentee, error) {
	mentee, err := m.menteeRepository.FindById(id)
	if err != nil {
		if err == helper.ErrMenteeNotFound {
			return nil, helper.ErrMenteeNotFound
		}

		return nil, helper.ErrInternalServerError
	}

	return mentee, nil
}

func (m menteeUseCase) FindByCourse(courseId string, pagination helper.Pagination) (*helper.Pagination, error) {
	mentees, totalRows, err := m.menteeRepository.FindByCourse(courseId, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, err
	}

	pagination.Result = mentees
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	return &pagination, nil
}

func (m menteeUseCase) Update(id string, menteeDomain *domain.Mentee) error {
	mentee, err := m.menteeRepository.FindById(id)
	if err != nil {
		return err
	}

	var ProfilePictureURL string

	if menteeDomain.ProfilePictureFile != nil {
		ctx := context.Background()

		if mentee.ProfilePicture != "" {
			if err := m.storage.DeleteObject(ctx, mentee.ProfilePicture); err != nil {
				return err
			}
		}

		ProfilePicture, err := menteeDomain.ProfilePictureFile.Open()
		if err != nil {
			return helper.ErrUnsupportedImageFile
		}

		defer ProfilePicture.Close()

		extension := filepath.Ext(menteeDomain.ProfilePictureFile.Filename)
		if extension != ".jpg" && extension != ".png" && extension != ".jpeg" {
			return helper.ErrUnsupportedImageFile
		}

		filename, _ := helper.GetFilename(menteeDomain.ProfilePictureFile.Filename)

		ProfilePictureURL, err = m.storage.UploadImage(ctx, filename, ProfilePicture)
		if err != nil {
			return err
		}
	}

	updatedMentee := domain.Mentee{
		Fullname:       menteeDomain.Fullname,
		Phone:          menteeDomain.Phone,
		BirthDate:      menteeDomain.BirthDate,
		Address:        menteeDomain.Address,
		ProfilePicture: ProfilePictureURL,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = m.menteeRepository.Update(id, &updatedMentee)
	if err != nil {
		if err == helper.ErrMenteeNotFound {
			return helper.ErrMenteeNotFound
		}

		return helper.ErrInternalServerError
	}

	return nil
}
