package domain

import (
	"mime/multipart"
	"time"
)

type Mentor struct {
	ID             string
	UserId         string
	Fullname       string
	Email          string
	Phone          string
	Role           string
	Jobs           string
	Gender         string
	BirthPlace     string
	BirthDate      time.Time
	Address        string
	ProfilePicture string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MentorAuth struct {
	Email    string
	Password string
}

type MentorRegister struct {
	Fullname string
	Email    string
	Password string
}

type MentorForgotPassword struct {
	Email string
}

type MentorUpdatePassword struct {
	UserID      string
	OldPassword string
	NewPassword string
}

type MentorUpdateProfile struct {
	UserID             string
	Fullname           string
	Email              string
	Phone              string
	Jobs               string
	Gender             string
	BirthPlace         string
	BirthDate          time.Time
	Address            string
	ProfilePictureFile *multipart.FileHeader
}

type MentorRepository interface {
	Create(mentorDomain *Mentor) error
	FindAll() (*[]Mentor, error)
	FindById(id string) (*Mentor, error)
	FindByIdUser(userId string) (*Mentor, error)
	Update(id string, mentorDomain *Mentor) error
}

type MentorUseCase interface {
	Register(mentorAuth *MentorRegister) error
	ForgotPassword(forgotPassword *MentorForgotPassword) error
	UpdatePassword(updatePassword *MentorUpdatePassword) error
	Login(mentorAuth *MentorAuth) (interface{}, error)
	FindAll() (*[]Mentor, error)
	FindById(id string) (*Mentor, error)
	Update(id string, updateMentor *MentorUpdateProfile) error
}
