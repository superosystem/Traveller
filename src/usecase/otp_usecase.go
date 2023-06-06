package usecase

import (
	"context"
	"fmt"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type otpUseCase struct {
	otpRepository  domain.OtpRepository
	userRepository domain.UserRepository
	mailer         *config.MailerConfig
}

func NewOTPUseCase(
	otpRepository domain.OtpRepository,
	userRepository domain.UserRepository,
	mailer *config.MailerConfig,
) domain.OtpUseCase {
	return otpUseCase{
		otpRepository:  otpRepository,
		userRepository: userRepository,
		mailer:         mailer,
	}
}

func (ou otpUseCase) SendOTP(otpDomain *domain.Otp) error {
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

	_ = ou.mailer.SendMail(user.Email, subject, message)

	return nil
}

func (ou otpUseCase) CheckOTP(otpDomain *domain.Otp) error {
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
