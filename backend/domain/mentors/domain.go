package mentors

import (
	"mime/multipart"
	"time"
)

type Domain struct {
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

type Repository interface {
	// Create repository create mentors
	Create(mentorDomain *Domain) error

	// FindAll repository find all mentors
	FindAll() (*[]Domain, error)

	// FindById repository find mentors by id
	FindById(id string) (*Domain, error)

	// FindByIdUser repository find mentors by id user
	FindByIdUser(userId string) (*Domain, error)

	// Update repository edit data mentors
	Update(id string, mentorDomain *Domain) error
}

type Usecase interface {
	// Register usecase mentors register
	Register(mentorAuth *MentorRegister) error

	// ForgotPassword usecase mentor verify forgot password
	ForgotPassword(forgotPassword *MentorForgotPassword) error

	// UpdatePassword usecase mentor to chnge password
	UpdatePassword(updatePassword *MentorUpdatePassword) error

	// Login usecase mentor login
	Login(mentorAuth *MentorAuth) (interface{}, error)

	// FindAll usecase find all mentors
	FindAll() (*[]Domain, error)

	// FindById usecase find by id mentors
	FindById(id string) (*Domain, error)

	// Update usecase edit data mentors
	Update(id string, updateMentor *MentorUpdateProfile) error
}
