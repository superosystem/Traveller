package usecase

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/superosystem/trainingsystem-backend/src/common/templates"
	"github.com/superosystem/trainingsystem-backend/src/domain"
	"html/template"
)

type certificateUseCase struct {
	menteeRepository domain.MenteeRepository
	courseRepository domain.CourseRepository
}

func NewCertificateUseCase(
	menteeRepository domain.MenteeRepository,
	courseRepository domain.CourseRepository,
) domain.CertificateUseCase {
	return certificateUseCase{
		menteeRepository: menteeRepository,
		courseRepository: courseRepository,
	}
}

func (cu certificateUseCase) GenerateCert(data *domain.Certificate) ([]byte, error) {
	mentee, err := cu.menteeRepository.FindById(data.MenteeId)
	if err != nil {
		return nil, err
	}

	course, err := cu.courseRepository.FindById(data.CourseId)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("").Parse(templates.Certificate)
	if err != nil {
		return nil, err
	}

	// prepare data certificate needs
	certDomain := map[string]string{
		"fullname": mentee.Fullname,
		"title":    course.Title,
	}

	// apply parsed HTML template data and keep the result in a buffer
	var w bytes.Buffer

	if err := tmpl.Execute(&w, certDomain); err != nil {
		return nil, err
	}

	// init a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewBuffer(w.Bytes()))

	page.EnableLocalFileAccess.Set(true)

	// add the page to generator
	pdfg.AddPage(page)

	// manipulate attribute
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA5)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Cover.Zoom.Set(1.2)
	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
