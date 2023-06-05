package domain

import (
	"mime/multipart"
	"time"

	"github.com/superosystem/trainingsystem-backend/src/helper"
)

type Mentee struct {
	ID                 string
	UserId             string
	Fullname           string
	Phone              string
	Role               string
	BirthDate          string
	Address            string
	ProfilePicture     string
	ProfilePictureFile *multipart.FileHeader
	User               User
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type MenteeAuth struct {
	Email    string
	Password string
}

type MenteeRegister struct {
	Fullname string
	Phone    string
	Email    string
	Password string
	OTP      string
}

type MenteeForgotPassword struct {
	Email            string
	Password         string
	RepeatedPassword string
	OTP              string
}

type MenteeRepository interface {
	// Create repository create mentee
	Create(menteeDomain *Mentee) error

	// FindAll repository find all mentees
	FindAll() (*[]Mentee, error)

	// FindById repository find mentee by id
	FindById(id string) (*Mentee, error)

	// FindByIdUser repository find mentee by id user
	FindByIdUser(userId string) (*Mentee, error)

	// repository find mentees by course
	FindByCourse(courseId string, limit int, offset int) (*[]Mentee, int, error)

	// repository count total mentees by course
	CountByCourse(courseId string) (int64, error)

	// Update repository edit data mentee
	Update(id string, menteeDomain *Mentee) error
}

type MenteeUsecase interface {
	// Register usecase mentee register
	Register(menteeAuth *MenteeAuth) error

	// VerifyRegister usecase verify register mentee
	VerifyRegister(menteeRegister *MenteeRegister) error

	// ForgotPassword usecase mentee verify forgot password
	ForgotPassword(forgotPassword *MenteeForgotPassword) error

	// Login usecase mentee login
	Login(menteeAuth *MenteeAuth) (interface{}, error)

	// FindAll usecase find all mentees
	FindAll() (*[]Mentee, error)

	// FindById usecase find by id mentee
	FindById(id string) (*Mentee, error)

	// usecase find mentee profile
	// MenteeProfile(menteeId string) (*Domain, error)

	// usecase find mentees by course
	FindByCourse(courseId string, pagination helper.Pagination) (*helper.Pagination, error)

	// Update usecase edit data mentee
	Update(id string, menteeDomain *Mentee) error
}
