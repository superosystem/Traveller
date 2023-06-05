package mentee_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/helper"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	menteeRepository     mocks.MenteeRepository
	menteeService        domain.MenteeUsecase
	otpRepository        mocks.OtpRepository
	userRepository       mocks.UserRepository
	jwtConfig            config.JWTConfig
	mailerConfig         config.MailerConfig
	storage              config.StorageConfig
	menteeDomain         domain.Mentee
	menteeAuth           domain.MenteeAuth
	menteeRegister       domain.MenteeRegister
	menteeForgotPassword domain.MenteeForgotPassword
	userDomain           domain.User

	pagination helper.Pagination
)

func TestMain(m *testing.M) {
	menteeService = usecase.NewMenteeUsecase(&menteeRepository, &userRepository, &otpRepository, &jwtConfig, &mailerConfig, &storage)

	userDomain = domain.User{
		ID:        uuid.NewString(),
		Email:     "test@gmail.com",
		Password:  "testtest",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	menteeDomain = domain.Mentee{
		ID:             uuid.NewString(),
		UserId:         userDomain.ID,
		Fullname:       "test",
		Phone:          "test",
		Role:           "mentee",
		Address:        "test",
		ProfilePicture: "test.com",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	menteeAuth = domain.MenteeAuth{
		Email:    userDomain.Email,
		Password: userDomain.Password,
	}

	menteeRegister = domain.MenteeRegister{
		Fullname: menteeDomain.Fullname,
		Phone:    menteeDomain.Phone,
		Email:    userDomain.Email,
		Password: userDomain.Password,
		OTP:      "0000",
	}

	menteeForgotPassword = domain.MenteeForgotPassword{
		Email:            userDomain.Email,
		Password:         userDomain.Password,
		RepeatedPassword: userDomain.Password,
		OTP:              "0000",
	}

	jwtConfig = config.JWTConfig{
		JWTSecret: "secret",
	}

	pagination = helper.Pagination{
		Limit: 10,
		Page:  1,
	}

	m.Run()
}

func TestRegister(t *testing.T) {
	t.Run("Test Register | Success register", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeAuth.Email).Return(nil, nil).Once()

		otpRepository.Mock.On("Save", mock.Anything, menteeAuth.Email, mock.Anything, helper.TIME_TO_LIVE).Return(nil).Once()

		err := menteeService.Register(&menteeAuth)

		assert.NoError(t, err)
	})

	t.Run("Test Register | Failed register | Invalid password length", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeAuth.Email).Return(nil, nil).Once()

		otpRepository.Mock.On("Save", mock.Anything, menteeAuth.Email, mock.Anything, helper.TIME_TO_LIVE).Return(nil).Once()

		menteeAuth.Password = "test"

		err := menteeService.Register(&menteeAuth)

		assert.Error(t, err)
	})

	t.Run("Test Register | Failed register | Error set otp", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeAuth.Email).Return(nil, nil).Once()

		otpRepository.Mock.On("Save", mock.Anything, menteeAuth.Email, mock.Anything, helper.TIME_TO_LIVE).Return(errors.New("error occured")).Once()

		err := menteeService.Register(&menteeAuth)

		assert.Error(t, err)
	})
}

func TestVerifyRegister(t *testing.T) {
	t.Run("Test VerifyRegister | Success register verified", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", userDomain.Email).Return(nil, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, menteeRegister.Email).Return(menteeRegister.OTP, nil).Once()

		userRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		menteeRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		err := menteeService.VerifyRegister(&menteeRegister)

		assert.NoError(t, err)
	})

	t.Run("Test VerifyRegister | Failed register | OTP not match", func(t *testing.T) {
		otpRepository.Mock.On("Get", mock.Anything, menteeRegister.Email).Return("9999", helper.ErrOTPNotMatch).Once()

		err := menteeService.VerifyRegister(&menteeRegister)

		assert.Error(t, err)
	})

	t.Run("Test VerifyRegister | Failed register verified | Error on create user", func(t *testing.T) {
		otpRepository.Mock.On("Get", mock.Anything, menteeRegister.Email).Return(menteeRegister.OTP, nil).Once()

		userRepository.Mock.On("Create", mock.Anything).Return(errors.New("error occurred")).Once()

		err := menteeService.VerifyRegister(&menteeRegister)

		assert.Error(t, err)
	})

	t.Run("Test VerifyRegister | Failed register verified | Error on create mentee", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", userDomain.Email).Return(nil, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, menteeRegister.Email).Return(menteeRegister.OTP, nil).Once()

		userRepository.Mock.On("Create", mock.Anything).Return(nil).Once()

		menteeRepository.Mock.On("Create", mock.Anything).Return(errors.New("error occurred"))

		err := menteeService.VerifyRegister(&menteeRegister)

		assert.Error(t, err)
	})
}

func TestForgotPassword(t *testing.T) {
	t.Run("Test Forgot Password | Success forgot password", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeForgotPassword.Email).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, menteeForgotPassword.Email).Return(menteeForgotPassword.OTP, nil).Once()

		userRepository.Mock.On("Update", userDomain.ID, mock.Anything).Return(nil).Once()

		err := menteeService.ForgotPassword(&menteeForgotPassword)

		assert.NoError(t, err)
	})

	t.Run("Test Forgot Password | Failed forgot password | User not found", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeForgotPassword.Email).Return(nil, helper.ErrUserNotFound).Once()

		err := menteeService.ForgotPassword(&menteeForgotPassword)

		assert.Error(t, err)
	})

	t.Run("Test Forgot Password | Failed forgot password | OTP expired", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeForgotPassword.Email).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, menteeForgotPassword.Email).Return("", helper.ErrOTPExpired).Once()

		err := menteeService.ForgotPassword(&menteeForgotPassword)

		assert.Error(t, err)
	})

	t.Run("Test Forgot Password | Failed forgot password | Error occurred", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeForgotPassword.Email).Return(&userDomain, nil).Once()

		otpRepository.Mock.On("Get", mock.Anything, menteeForgotPassword.Email).Return(menteeForgotPassword.OTP, nil).Once()

		userRepository.Mock.On("Update", userDomain.ID, mock.Anything).Return(errors.New("error occurred"))

		err := menteeService.ForgotPassword(&menteeForgotPassword)

		assert.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Test Login | Failed login | Invalid password length", func(t *testing.T) {
		userRepository.Mock.On("FindByEmail", menteeAuth.Email).Return(&userDomain, nil).Once()

		menteeRepository.Mock.On("FindByIdUser", userDomain.ID).Return(&menteeDomain, nil).Once()

		result, err := menteeService.Login(&menteeAuth)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("Test Find All | Success find all mentees", func(t *testing.T) {
		menteeRepository.Mock.On("FindAll").Return(&[]domain.Mentee{menteeDomain}, nil).Once()

		results, err := menteeService.FindAll()

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find All | Failed find all mentees | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindAll").Return(nil, helper.ErrMenteeNotFound).Once()

		results, err := menteeService.FindAll()

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindById(t *testing.T) {
	t.Run("Test Find By Id | Success find by id", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		results, err := menteeService.FindById(menteeDomain.ID)

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find By Id | Failed find by id | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(nil, helper.ErrMenteeNotFound).Once()

		results, err := menteeService.FindById(menteeDomain.ID)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestFindByCourse(t *testing.T) {
	t.Run("Test Find By Course | Success find by course", func(t *testing.T) {
		menteeRepository.Mock.On("FindByCourse", "test", pagination.GetLimit(), pagination.GetOffset()).Return(&[]domain.Mentee{menteeDomain}, 1, nil).Once()

		results, err := menteeService.FindByCourse("test", pagination)

		assert.NoError(t, err)
		assert.NotEmpty(t, results)
	})

	t.Run("Test Find By Course | Failed find by course | Error occurred", func(t *testing.T) {
		menteeRepository.Mock.On("FindByCourse", "test", pagination.GetLimit(), pagination.GetOffset()).Return(&[]domain.Mentee{}, 0, errors.New("error occurred")).Once()

		results, err := menteeService.FindByCourse("test", pagination)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Test Update | Success update", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		menteeRepository.Mock.On("Update", menteeDomain.ID, mock.Anything).Return(nil).Once()

		err := menteeService.Update(menteeDomain.ID, &menteeDomain)

		assert.NoError(t, err)
	})

	t.Run("Test Update | Failed update | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(nil, helper.ErrMenteeNotFound).Once()

		err := menteeService.Update(menteeDomain.ID, &menteeDomain)

		assert.Error(t, err)
	})

	t.Run("Test Update | Failed update | Mentee not found", func(t *testing.T) {
		menteeRepository.Mock.On("FindById", menteeDomain.ID).Return(&menteeDomain, nil).Once()

		menteeRepository.Mock.On("Update", menteeDomain.ID, mock.Anything).Return(errors.New("error occurred")).Once()

		err := menteeService.Update(menteeDomain.ID, &menteeDomain)

		assert.Error(t, err)
	})
}
