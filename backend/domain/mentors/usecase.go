package mentors

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/superosystem/TrainingSystem/backend/domain/users"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type mentorUsecase struct {
	mentorsRepository Repository
	userRepository    users.Repository
	jwtConfig         *helper.JWTConfig
	storage           *helper.StorageConfig
	mailerConfig      *helper.MailerConfig
}

func NewMentorUsecase(mentorsRepository Repository, userRepository users.Repository, jwtConfig *helper.JWTConfig, storage *helper.StorageConfig, mailerConfig *helper.MailerConfig) Usecase {
	return mentorUsecase{
		mentorsRepository: mentorsRepository,
		userRepository:    userRepository,
		jwtConfig:         jwtConfig,
		storage:           storage,
		mailerConfig:      mailerConfig,
	}
}

func (m mentorUsecase) Register(mentorDomain *MentorRegister) error {
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

	user := users.Domain{
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

	mentor := Domain{
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

func (m mentorUsecase) Login(mentorAuth *MentorAuth) (interface{}, error) {
	if len(mentorAuth.Password) < 6 {
		return nil, helper.ErrPasswordLengthInvalid
	}

	var err error

	var user *users.Domain
	user, err = m.userRepository.FindByEmail(mentorAuth.Email)

	if err != nil {
		return nil, err
	}

	ok := helper.ComparePassword(user.Password, mentorAuth.Password)
	if !ok {
		return nil, helper.ErrUserNotFound
	}

	var mentor *Domain
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

func (m mentorUsecase) UpdatePassword(updatePassword *MentorUpdatePassword) error {

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

	updatedUser := users.Domain{
		Password: hashPassword,
	}

	err = m.userRepository.Update(oldPassword.ID, &updatedUser)

	if err != nil {
		return err
	}

	return nil

}

func (m mentorUsecase) FindAll() (*[]Domain, error) {
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

func (m mentorUsecase) FindById(id string) (*Domain, error) {

	mentor, err := m.mentorsRepository.FindById(id)
	if err != nil {
		if err == helper.ErrMentorNotFound {
			return nil, helper.ErrMentorNotFound
		}

		return nil, helper.ErrInternalServerError
	}

	return mentor, nil
}

func (m mentorUsecase) Update(id string, updateMentor *MentorUpdateProfile) error {
	_, err := m.userRepository.FindById(updateMentor.UserID)

	if err != nil {
		return err
	}
	user := users.Domain{
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

	updatedMentor := Domain{

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

func (m mentorUsecase) ForgotPassword(forgotPassword *MentorForgotPassword) error {
	var err error

	user, err := m.userRepository.FindByEmail(forgotPassword.Email)

	if err != nil {
		return err
	}

	randomPassword := helper.GenerateOTP(8)
	hashPassword := helper.HashPassword(randomPassword)

	updatedUser := users.Domain{
		Password: hashPassword,
	}

	err = m.userRepository.Update(user.ID, &updatedUser)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("Password baru anda: %s", randomPassword)
	subject := "Reset Password T"

	_ = m.mailerConfig.SendMail(user.Email, subject, message)

	return nil
}
