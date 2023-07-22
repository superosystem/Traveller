package otp

import (
	"context"
	"fmt"

	"github.com/superosystem/TrainingSystem/backend/constants"
	"github.com/superosystem/TrainingSystem/backend/domain/users"
	"github.com/superosystem/TrainingSystem/backend/helper"
)

type otpUsecase struct {
	otpRepository  Repository
	userRepository users.Repository
	mailerconfig   *helper.MailerConfig
}

func NewOTPUsecase(
	otpRepository Repository,
	userRepository users.Repository,
	mailerconfig *helper.MailerConfig,
) Usecase {
	return otpUsecase{
		otpRepository:  otpRepository,
		userRepository: userRepository,
		mailerconfig:   mailerconfig,
	}
}

func (ou otpUsecase) SendOTP(otpDomain *Domain) error {
	var err error

	var user *users.Domain
	user, err = ou.userRepository.FindByEmail(otpDomain.Key)

	if err != nil {
		return err
	}

	ctx := context.Background()
	newOTP := helper.GenerateOTP(4)

	err = ou.otpRepository.Save(ctx, user.Email, newOTP, constants.TIME_TO_LIVE)

	if err != nil {
		return err
	}

	subject := "Verification Code Eduworld"
	message := fmt.Sprintf("OTP: %s", newOTP)

	_ = ou.mailerconfig.SendMail(user.Email, subject, message)

	return nil
}

func (ou otpUsecase) CheckOTP(otpDomain *Domain) error {
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
