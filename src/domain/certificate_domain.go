package domain

type Certificate struct {
	MenteeId string
	CourseId string
}

// PDFService represents the interface of a pdf generation service
type CertificateUsecase interface {
	GenerateCert(data *Certificate) ([]byte, error)
}
