package domain

import (
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"mime/multipart"
	"time"
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
	Create(menteeDomain *Mentee) error
	FindAll() (*[]Mentee, error)
	FindById(id string) (*Mentee, error)
	FindByIdUser(userId string) (*Mentee, error)
	FindByCourse(courseId string, limit int, offset int) (*[]Mentee, int, error)
	CountByCourse(courseId string) (int64, error)
	Update(id string, menteeDomain *Mentee) error
}

type MenteeUseCase interface {
	Register(menteeAuth *MenteeAuth) error
	VerifyRegister(menteeRegister *MenteeRegister) error
	ForgotPassword(forgotPassword *MenteeForgotPassword) error
	Login(menteeAuth *MenteeAuth) (interface{}, error)
	FindAll() (*[]Mentee, error)
	FindById(id string) (*Mentee, error)
	FindByCourse(courseId string, pagination helper.Pagination) (*helper.Pagination, error)
	Update(id string, menteeDomain *Mentee) error
}
