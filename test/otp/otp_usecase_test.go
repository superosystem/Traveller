package otp_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	otpRepository  mocks.OtpRepository
	userRepository mocks.UserRepository
	mailerConfig   config.MailerConfig
	otpService     domain.OtpUseCase
	otpDomain      domain.Otp
	userDomain     domain.User
)

func TestMain(m *testing.M) {
	otpService = usecase.NewOTPUseCase(&otpRepository, &userRepository, &mailerConfig)

	otpDomain = domain.Otp{
		Key:   "test@gmail.com",
		Value: "0000",
	}

	userDomain = domain.User{
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

		otpRepository.Mock.On("Save", mock.Anything, otpDomain.Key, mock.Anything, helper.TIME_TO_LIVE).Return(nil).Once()

		err := otpService.SendOTP(&otpDomain)

		assert.NoError(t, err)
	})

	t.Run("Test SendOTP | Failed send otp | User not found", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&domain.User{}, helper.ErrUserNotFound).Once()

		err := otpService.SendOTP(&otpDomain)

		assert.Error(t, err)
	})

	t.Run("Test SendOTP | Failed send otp | Error occurred", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Save", mock.Anything, otpDomain.Key, mock.Anything, helper.TIME_TO_LIVE).Return(errors.New("error occurred")).Once()

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
		userRepository.Mock.On("FindByEmail", otpDomain.Key).Return(&domain.User{}, helper.ErrUserNotFound).Once()

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
