package certificates

type Domain struct {
	MenteeId string
	CourseId string
}

// PDFService represents the interface of a pdf generation service
type Usecase interface {
	GenerateCert(data *Domain) ([]byte, error)
}
