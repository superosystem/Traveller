package mentor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	mocks "github.com/superosystem/trainingsystem-backend/src/domain/mocks/repository"
	"github.com/superosystem/trainingsystem-backend/src/usecase"
)

var (
	mentorRepository mocks.MentorRepository
	mentorService    domain.MentorUseCase
	userRepository   mocks.UserRepository
	jwtConfig        config.JWTConfig
	storageClient    config.StorageConfig
	mailerConfig     config.MailerConfig
	mentorDomain     domain.Mentor
	mentorRegister   domain.MentorRegister
	_                domain.MentorUpdateProfile
)

func TestMain(m *testing.M) {
	mentorRepository = mocks.MentorRepository{Mock: mock.Mock{}}
	userRepository = mocks.UserRepository{Mock: mock.Mock{}}
	jwtConfig = config.JWTConfig{JWTSecret: "secret"}
	storageClient = config.StorageConfig{}

	mentorService = usecase.NewMentorUseCase(&mentorRepository, &userRepository, &jwtConfig, &storageClient, &mailerConfig)

	// birth date
	birthDate := time.Date(2021, 8, 11, 0, 0, 0, 0, time.Local)

	mentorDomain = domain.Mentor{
		ID:             "MID1",
		UserId:         "UID1",
		Fullname:       "Mentors Test",
		Phone:          "0857654378",
		Role:           "mentor",
		BirthDate:      birthDate,
		Address:        "Jl Ahmad Yani",
		ProfilePicture: "https://example.com/to/bucket",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// mentorAuth = mentors.MentorAuth{
	// 	Email:    "mentor@gmail.com",
	// 	Password: "hashedpassword",
	// }

	mentorRegister = domain.MentorRegister{
		Fullname: "Mentor Test",
		Email:    "mentor@gmail.com",
		Password: "hashedpassword",
	}

	m.Run()
}

func TestRegister(t *testing.T) {
	t.Run("Register | Success", func(t *testing.T) {

		userRepository.Mock.On("FindByEmail", mentorRegister.Email).Return(nil, nil)

		userRepository.Mock.On("Create", mock.Anything).Return(nil)

		mentorRepository.Mock.On("Create", mock.Anything).Return(nil)

		err := mentorService.Register(&mentorRegister)

		assert.NoError(t, err)
	})
}

// TODO: Create test Login

// func TestLogin(t *testing.T) {
// 	t.Run("Login | Success", func(t *testing.T) {

// 		userDomain = users.Domain{
// 			ID:        "UID1",
// 			Email:     "mentor@gmail.com",
// 			Password:  "hashedpassword",
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}

// 		mentorAuth = mentors.MentorAuth{
// 			Email:    "mentor@gmail.com",
// 			Password: "hashedpassword",
// 		}

// 		userRepository.Mock.On("FindByEmail", mentorAuth.Email).Return(&userDomain, nil)

// 		token, err := mentorService.Login(&mentorAuth)

// 		assert.Error(t, err)
// 		assert.NotNil(t, &token)
// 	})
// }

func TestFindById(t *testing.T) {
	t.Run("Find By ID | Valid", func(t *testing.T) {
		mentorRepository.On("FindById", "MID1").Return(&mentorDomain, nil).Once()

		result, err := mentorService.FindById("MID1")
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Find By ID | InValid", func(t *testing.T) {
		mentorRepository.On("FindById", "-1").Return(&mentorDomain, nil).Once()

		result, err := mentorService.FindById("-1")
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("Find All | Valid", func(t *testing.T) {
		mentorRepository.On("FindAll").Return(&[]domain.Mentor{mentorDomain}, nil).Once()

		result, err := mentorRepository.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, 1, len(*result))
	})

	t.Run("Find All | InValid", func(t *testing.T) {
		mentorRepository.On("FindAll").Return(&[]domain.Mentor{}, nil).Once()

		result, err := mentorRepository.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(*result))
	})
}
