package otp

import "github.com/superosystem/TrainingSystem/backend/domain/otp"

type OTPCode struct {
	Key string `json:"email"`
}

func FromDomain(otpDomain *otp.Domain) *OTPCode {
	return &OTPCode{
		Key: otpDomain.Key,
	}
}

func (rec *OTPCode) ToDomain() *otp.Domain {
	return &otp.Domain{
		Key: rec.Key,
	}
}
