package manage_mentees

type Usecase interface {
	// usecase delete access course mentee
	DeleteAccess(menteeId string, courseId string) error
}
