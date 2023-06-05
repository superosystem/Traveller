package domain

type ManageMenteeUsecase interface {
	// usecase delete access course mentee
	DeleteAccess(menteeId string, courseId string) error
}
