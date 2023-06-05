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
	// Create repository create mentors
	Create(mentorDomain *Mentor) error

	// FindAll repository find all mentors
	FindAll() (*[]Mentor, error)

	// FindById repository find mentors by id
	FindById(id string) (*Mentor, error)

	// FindByIdUser repository find mentors by id user
	FindByIdUser(userId string) (*Mentor, error)

	// Update repository edit data mentors
	Update(id string, mentorDomain *Mentor) error
}

type MentorUsecase interface {
	// Register usecase mentors register
	Register(mentorAuth *MentorRegister) error

	// ForgotPassword usecase mentor verify forgot password
	ForgotPassword(forgotPassword *MentorForgotPassword) error

	// UpdatePassword usecase mentor to chnge password
	UpdatePassword(updatePassword *MentorUpdatePassword) error

	// Login usecase mentor login
	Login(mentorAuth *MentorAuth) (interface{}, error)

	// FindAll usecase find all mentors
	FindAll() (*[]Mentor, error)

	// FindById usecase find by id mentors
	FindById(id string) (*Mentor, error)

	// Update usecase edit data mentors
	Update(id string, updateMentor *MentorUpdateProfile) error
}
