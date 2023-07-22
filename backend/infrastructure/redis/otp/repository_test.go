package otp

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"
	"github.com/superosystem/TrainingSystem/backend/constants"
	otpDomain "github.com/superosystem/TrainingSystem/backend/domain/otp"
)

type suiteOtp struct {
	suite.Suite
	mock          redismock.ClientMock
	otpRepository otpDomain.Repository
}

func (s *suiteOtp) SetupSuite() {
	db, mock := redismock.NewClientMock()

	s.mock = mock

	s.otpRepository = NewRedisRepository(db)
}

func (s *suiteOtp) TestSave() {
	key := "mentee@gmail.com"
	value := "7339"

	s.mock.Regexp().ExpectSet(key, value, constants.TIME_TO_LIVE).
		SetVal(value)

	ctx := context.TODO()

	err := s.otpRepository.Save(ctx, key, value, constants.TIME_TO_LIVE)

	s.NoError(err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.Error(err)
	}
}

func (s *suiteOtp) TestGet() {
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
