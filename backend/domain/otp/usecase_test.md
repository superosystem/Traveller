package otp_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/TrainingSystem/backend/constants"
	"github.com/superosystem/TrainingSystem/backend/domain/otp"
	_otpMock "github.com/superosystem/TrainingSystem/backend/domain/otp/mocks"
	"github.com/superosystem/TrainingSystem/backend/domain/users"
	_userMock "github.com/superosystem/TrainingSystem/backend/domain/users/mocks"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

var (
	otpRepository  _otpMock.Repository
	userRepository _userMock.UserRepositoryMock
	mailerConfig   helper.MailerConfig

	otpService otp.Usecase

	otpDomain  otp.Domain
	userDomain users.Domain
)

func TestMain(m *testing.M) {
	otpService = otp.NewOTPUsecase(&otpRepository, &userRepository, &mailerConfig)

	otpDomain = otp.Domain{
		Key:   "test@gmail.com",
		Value: "0000",
	}

	userDomain = users.Domain{
		ID:        uuid.NewString(),
		Email:     otpDomain.Key,
		Password:  "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.Run()
}

func TestSendOTP(t *testing.T) {
	t.Run("Test SendOTP | Success send otp", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Save", mock.Anything, otpDomain.Key, mock.Anything, constants.TIME_TO_LIVE).Return(nil).Once()

		err := otpService.SendOTP(&otpDomain)

		assert.NoError(t, err)
	})

	t.Run("Test SendOTP | Failed send otp | User not found", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&users.Domain{}, helper.ErrUserNotFound).Once()

		err := otpService.SendOTP(&otpDomain)

		assert.Error(t, err)
	})

	t.Run("Test SendOTP | Failed send otp | Error occurred", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Save", mock.Anything, otpDomain.Key, mock.Anything, constants.TIME_TO_LIVE).Return(errors.New("error occurred")).Once()

		err := otpService.SendOTP(&otpDomain)

		assert.Error(t, err)
	})
}

func TestCheckOTP(t *testing.T) {
	t.Run("Test CheckOTP | Success check otp", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, otpDomain.Key).Return(otpDomain.Value, nil).Once()

		err := otpService.CheckOTP(&otpDomain)

		assert.NoError(t, err)
	})

	t.Run("Test CheckOTP | Failed check otp | User Not Found", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&users.Domain{}, helper.ErrUserNotFound).Once()

		err := otpService.CheckOTP(&otpDomain)

		assert.Error(t, err)
	})

	t.Run("Test CheckOTP | Failed check otp | OTP expired", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, otpDomain.Key).Return("", helper.ErrOTPExpired).Once()

		err := otpService.CheckOTP(&otpDomain)

		assert.Error(t, err)
	})

	t.Run("Test CheckOTP | Failed check otp | OTP not match", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, otpDomain.Key).Return("9999", helper.ErrOTPNotMatch).Once()

		err := otpService.CheckOTP(&otpDomain)

		assert.Error(t, err)
	})
}
