package usecase

import (
	"context"
	"fmt"

	"github.com/superosystem/trainingsystem-backend/src/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type otpUsecase struct {
	otpRepository  domain.OtpRepository
	userRepository domain.UserRepository
	mailerconfig   *config.MailerConfig
}

func NewOTPUsecase(
	otpRepository domain.OtpRepository,
	userRepository domain.UserRepository,
	mailerconfig *config.MailerConfig,
) domain.OtpUsecase {
	return otpUsecase{
		otpRepository:  otpRepository,
		userRepository: userRepository,
		mailerconfig:   mailerconfig,
	}
}

func (ou otpUsecase) SendOTP(otpDomain *domain.Otp) error {
	var err error

	var user *domain.User
	user, err = ou.userRepository.FindByEmail(otpDomain.Key)

	if err != nil {
		return err
	}

	ctx := context.Background()
	newOTP := helper.GenerateOTP(4)

	err = ou.otpRepository.Save(ctx, user.Email, newOTP, helper.TIME_TO_LIVE)

	if err != nil {
		return err
	}

	subject := "Verification Code Training System"
	message := fmt.Sprintf("OTP: %s", newOTP)

	_ = ou.mailerconfig.SendMail(user.Email, subject, message)

	return nil
}

func (ou otpUsecase) CheckOTP(otpDomain *domain.Otp) error {
	if _, err := ou.userRepository.FindByEmail(otpDomain.Key); err != nil {
		return err
	}

	ctx := context.Background()

	result, err := ou.otpRepository.Get(ctx, otpDomain.Key)

	if err != nil {
		return err
	}

	if result != otpDomain.Value {
		return helper.ErrOTPNotMatch
	}

	return nil
}
