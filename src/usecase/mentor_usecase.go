package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/superosystem/trainingsystem-backend/src/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type mentorUsecase struct {
	mentorsRepository domain.MentorRepository
	userRepository    domain.UserRepository
	jwtConfig         *config.JWTConfig
	storage           *config.StorageConfig
	mailerConfig      *config.MailerConfig
}

func NewMentorUsecase(
	mentorsRepository domain.MentorRepository,
	userRepository domain.UserRepository,
	jwtConfig *config.JWTConfig,
	storage *config.StorageConfig,
	mailerConfig *config.MailerConfig,
) domain.MentorUsecase {
	return mentorUsecase{
		mentorsRepository: mentorsRepository,
		userRepository:    userRepository,
		jwtConfig:         jwtConfig,
		storage:           storage,
		mailerConfig:      mailerConfig,
	}
}

func (m mentorUsecase) Register(mentorDomain *domain.MentorRegister) error {
	var err error

	if len(mentorDomain.Password) < 6 {
		return helper.ErrPasswordLengthInvalid
	}

	email, _ := m.userRepository.FindByEmail(mentorDomain.Email)

	if email != nil {
		return helper.ErrEmailAlreadyExist
	}

	userId := uuid.NewString()
	hashedPassword := helper.HashPassword(mentorDomain.Password)

	user := domain.User{
		ID:        userId,
		Email:     mentorDomain.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = m.userRepository.Create(&user)

	if err != nil {
		return err
	}

	mentorId := uuid.NewString()

	mentor := domain.Mentor{
		ID:        mentorId,
		UserId:    userId,
		Fullname:  mentorDomain.Fullname,
		Role:      "mentor",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = m.mentorsRepository.Create(&mentor)

	if err != nil {
		return err
	}

	return nil
}

func (m mentorUsecase) Login(mentorAuth *domain.MentorAuth) (interface{}, error) {
	if len(mentorAuth.Password) < 6 {
		return nil, helper.ErrPasswordLengthInvalid
	}

	var err error

	var user *domain.User
	user, err = m.userRepository.FindByEmail(mentorAuth.Email)

	if err != nil {
		return nil, err
	}

	ok := helper.ComparePassword(user.Password, mentorAuth.Password)
	if !ok {
		return nil, helper.ErrUserNotFound
	}

	var mentor *domain.Mentor
	mentor, err = m.mentorsRepository.FindByIdUser(user.ID)

	if err != nil {
		return nil, err
	}

	var token string
	exp := time.Now().Add(6 * time.Hour)

	token, err = m.jwtConfig.GenerateToken(user.ID, mentor.ID, mentor.Role, exp)

	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"token":   token,
		"expires": exp,
	}

	return data, nil
}

func (m mentorUsecase) UpdatePassword(updatePassword *domain.MentorUpdatePassword) error {

	if len(updatePassword.NewPassword) < 8 {
		return helper.ErrPasswordLengthInvalid
	}

	oldPassword, err := m.userRepository.FindById(updatePassword.UserID)

	if err != nil {
		return helper.ErrUserNotFound
	}

	ok := helper.ComparePassword(oldPassword.Password, updatePassword.OldPassword)
	if !ok {
		return helper.ErrPasswordNotMatch
	}

	hashPassword := helper.HashPassword(updatePassword.NewPassword)

	updatedUser := domain.User{
		Password: hashPassword,
	}

	err = m.userRepository.Update(oldPassword.ID, &updatedUser)

	if err != nil {
		return err
	}

	return nil

}

func (m mentorUsecase) FindAll() (*[]domain.Mentor, error) {
	var err error

	mentor, err := m.mentorsRepository.FindAll()

	if err != nil {
		if err == helper.ErrMentorNotFound {
			return nil, helper.ErrMentorNotFound
		}

		return nil, helper.ErrInternalServerError
	}

	return mentor, nil
}

func (m mentorUsecase) FindById(id string) (*domain.Mentor, error) {

	mentor, err := m.mentorsRepository.FindById(id)
	if err != nil {
		if err == helper.ErrMentorNotFound {
			return nil, helper.ErrMentorNotFound
		}

		return nil, helper.ErrInternalServerError
	}

	return mentor, nil
}

func (m mentorUsecase) Update(id string, updateMentor *domain.MentorUpdateProfile) error {
	_, err := m.userRepository.FindById(updateMentor.UserID)

	if err != nil {
		return err
	}
	user := domain.User{
		ID:        updateMentor.UserID,
		Email:     updateMentor.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = m.userRepository.Update(user.ID, &user)
	if err != nil {
		return err
	}

	mentor, err := m.mentorsRepository.FindById(id)

	if err != nil {
		return err
	}

	var ProfilePictureURL string

	if updateMentor.ProfilePictureFile != nil {
		ctx := context.Background()

		if mentor.ProfilePicture != "" {
			if err := m.storage.DeleteObject(ctx, mentor.ProfilePicture); err != nil {
				return err
			}
		}

		ProfilePicture, err := updateMentor.ProfilePictureFile.Open()

		if err != nil {
			return err
		}

		defer ProfilePicture.Close()

		extension := filepath.Ext(updateMentor.ProfilePictureFile.Filename)

		if extension != ".jpg" && extension != ".png" && extension != ".jpeg" {
			return helper.ErrUnsupportedImageFile
		}

		filename, _ := helper.GetFilename(updateMentor.ProfilePictureFile.Filename)

		ProfilePictureURL, err = m.storage.UploadImage(ctx, filename, ProfilePicture)

		if err != nil {
			return err
		}
	}

	updatedMentor := domain.Mentor{
		Fullname:       updateMentor.Fullname,
		Phone:          updateMentor.Phone,
		Jobs:           updateMentor.Jobs,
		Gender:         updateMentor.Gender,
		BirthPlace:     updateMentor.BirthPlace,
		BirthDate:      updateMentor.BirthDate,
		Address:        updateMentor.Address,
		ProfilePicture: ProfilePictureURL,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = m.mentorsRepository.Update(id, &updatedMentor)
	if err != nil {
		if err == helper.ErrMentorNotFound {
			return helper.ErrMentorNotFound
		}

		return helper.ErrInternalServerError
	}

	return nil
}

func (m mentorUsecase) ForgotPassword(forgotPassword *domain.MentorForgotPassword) error {
	var err error

	user, err := m.userRepository.FindByEmail(forgotPassword.Email)

	if err != nil {
		return err
	}

	randomPassword := helper.GenerateOTP(8)
	hashPassword := helper.HashPassword(randomPassword)

	updatedUser := domain.User{
		Password: hashPassword,
	}

	err = m.userRepository.Update(user.ID, &updatedUser)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("Password baru anda: %s", randomPassword)
	subject := "Reset Password Training System"

	_ = m.mailerConfig.SendMail(user.Email, subject, message)

	return nil
}
