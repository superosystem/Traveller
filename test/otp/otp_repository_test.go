package otp_test

import (
	"context"
	"github.com/superosystem/trainingsystem-backend/src/repository"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"github.com/superosystem/trainingsystem-backend/src/domain"
)

type suiteOtp struct {
	suite.Suite
	mock          redismock.ClientMock
	otpRepository domain.OtpRepository
}

func (s *suiteOtp) SetupSuite() {
	db, mock := redismock.NewClientMock()

	s.mock = mock

	s.otpRepository = repository.NewOtpRepository(db)
}

func (s *suiteOtp) TestSaveOTP() {
	key := "mentee@gmail.com"
	value := "7339"

	s.mock.Regexp().ExpectSet(key, value, helper.TIME_TO_LIVE).
		SetVal(value)

	ctx := context.TODO()

	err := s.otpRepository.Save(ctx, key, value, helper.TIME_TO_LIVE)

	s.NoError(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Error(err)
	}
}

func (s *suiteOtp) TestGetOTP() {
	key := "mentee@gmail.com"
	value := "7339"

	s.mock.ExpectGet(key).
		SetVal(value)

	ctx := context.TODO()

	result, err := s.otpRepository.Get(ctx, key)

	s.Nil(err)
	s.NotNil(result)

	s.Equal(value, result)
}

func TestSuiteOTP(t *testing.T) {
	suite.Run(t, new(suiteOtp))
}
