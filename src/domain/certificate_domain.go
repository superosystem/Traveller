package domain

type Certificate struct {
	MenteeId string
	CourseId string
}

// CertificateUseCase Generate Certificate to PDF
type CertificateUseCase interface {
	GenerateCert(data *Certificate) ([]byte, error)
}
