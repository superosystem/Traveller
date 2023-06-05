package entities

import "github.com/superosystem/trainingsystem-backend/src/domain"

type OTPCode struct {
	Key string `json:"email"`
}

func FromOtpDomain(otpDomain *domain.Otp) *OTPCode {
	return &OTPCode{
		Key: otpDomain.Key,
	}
}

func (rec *OTPCode) ToOtpDomain() *domain.Otp {
	return &domain.Otp{
		Key: rec.Key,
	}
}
