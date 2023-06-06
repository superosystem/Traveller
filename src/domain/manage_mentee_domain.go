package domain

type ManageMenteeUseCase interface {
	DeleteAccess(menteeId string, courseId string) error
}
