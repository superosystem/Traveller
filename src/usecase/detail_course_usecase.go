package usecase

import "github.com/superosystem/trainingsystem-backend/src/domain"

type detailCourseUsecase struct {
	menteeRepository           domain.MenteeRepository
	courseRepository           domain.CourseRepository
	moduleRepository           domain.ModuleRepository
	materialRepository         domain.MaterialRepository
	menteeProgressRepository   domain.MenteeProgressRepository
	assignmentsRepository      domain.AssignmentRepository
	menteeAssignmentRepository domain.MenteeAssignmentRepository
	menteeCourse               domain.MenteeCourseRepository
}

func NewDetailCourseUsecase(
	menteeRepository domain.MenteeRepository,
	courseRepository domain.CourseRepository,
	moduleRepository domain.ModuleRepository,
	materialRepository domain.MaterialRepository,
	menteeProgressRepository domain.MenteeProgressRepository,
	assignmentsRepository domain.AssignmentRepository,
	menteeAssignmentRepository domain.MenteeAssignmentRepository,
	menteeCourse domain.MenteeCourseRepository,
) domain.DetailCourseUsecase {
	return detailCourseUsecase{
		menteeRepository:           menteeRepository,
		courseRepository:           courseRepository,
		moduleRepository:           moduleRepository,
		materialRepository:         materialRepository,
		menteeProgressRepository:   menteeProgressRepository,
		assignmentsRepository:      assignmentsRepository,
		menteeAssignmentRepository: menteeAssignmentRepository,
		menteeCourse:               menteeCourse,
	}
}

func (dc detailCourseUsecase) DetailCourse(courseId string) (*domain.DetailCourse, error) {
	course, err := dc.courseRepository.FindById(courseId)

	if err != nil {
		return nil, err
	}

	modules, _ := dc.moduleRepository.FindByCourse(courseId)

	moduleIds := []string{}

	for _, module := range modules {
		moduleIds = append(moduleIds, module.ID)
	}

	assignment, _ := dc.assignmentsRepository.FindByCourseId(courseId)

	assignmentDomain := domain.DetailAssignment{}

	if assignment != nil {
		assignmentDomain = domain.DetailAssignment{
			ID:          assignment.ID,
			CourseId:    assignment.CourseId,
			Title:       assignment.Title,
			Description: assignment.Description,
			CreatedAt:   assignment.CreatedAt,
			UpdatedAt:   assignment.UpdatedAt,
		}
	}

	materials, _ := dc.materialRepository.FindByModule(moduleIds)

	materialDomain := make([]domain.DetailMaterial, len(materials))

	for i, material := range materials {
		materialDomain[i].MaterialId = material.ID
		materialDomain[i].ModuleId = material.ModuleId
		materialDomain[i].Title = material.Title
		materialDomain[i].URL = material.URL
		materialDomain[i].Description = material.Description
		materialDomain[i].CreatedAt = material.CreatedAt
		materialDomain[i].UpdatedAt = material.UpdatedAt
	}

	moduleDomain := make([]domain.DetailModule, len(modules))

	for i, module := range modules {
		moduleDomain[i].ModuleId = module.ID
		moduleDomain[i].CourseId = module.CourseId
		moduleDomain[i].Title = module.Title
		moduleDomain[i].Description = module.Description
		moduleDomain[i].CreatedAt = module.CreatedAt
		moduleDomain[i].UpdatedAt = module.UpdatedAt
	}

	for i, module := range moduleDomain {
		for _, material := range materialDomain {
			if module.ModuleId == material.ModuleId {
				moduleDomain[i].Materials = append(moduleDomain[i].Materials, material)
			}
		}
	}

	courseDomain := domain.DetailCourse{
		CourseId:     course.ID,
		CategoryId:   course.CategoryId,
		MentorId:     course.MentorId,
		Mentor:       course.Mentor.Fullname,
		Category:     course.Category.Name,
		Title:        course.Title,
		Description:  course.Description,
		Thumbnail:    course.Thumbnail,
		TotalReviews: course.TotalReviews,
		Rating:       course.Rating,
		Modules:      moduleDomain,
		Assignment:   assignmentDomain,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}

	return &courseDomain, nil
}

func (dc detailCourseUsecase) DetailCourseEnrolled(menteeId string, courseId string) (*domain.DetailCourse, error) {
	menteeCourse, err := dc.menteeCourse.CheckEnrollment(menteeId, courseId)

	if err != nil {
		return nil, err
	}

	course, err := dc.courseRepository.FindById(courseId)

	if err != nil {
		return nil, err
	}

	if _, err := dc.menteeRepository.FindById(menteeId); err != nil {
		return nil, err
	}

	assignment, _ := dc.assignmentsRepository.FindByCourseId(courseId)

	menteeAssignment, _ := dc.menteeAssignmentRepository.FindByCourse(menteeId, courseId)

	isCompletingAssignment := menteeAssignment != nil

	assignmentDomain := domain.DetailAssignment{}

	if assignment != nil {
		assignmentDomain = domain.DetailAssignment{
			ID:          assignment.ID,
			CourseId:    assignment.CourseId,
			Title:       assignment.Title,
			Description: assignment.Description,
			Completed:   isCompletingAssignment,
			CreatedAt:   assignment.CreatedAt,
			UpdatedAt:   assignment.UpdatedAt,
		}
	}

	modules, _ := dc.moduleRepository.FindByCourse(courseId)

	modulesIds := []string{}

	for _, module := range modules {
		modulesIds = append(modulesIds, module.ID)
	}

	materials, _ := dc.materialRepository.FindByModule(modulesIds)

	materialDomain := make([]domain.DetailMaterial, len(materials))

	for i, material := range materials {
		materialDomain[i].MaterialId = material.ID
		materialDomain[i].ModuleId = material.ModuleId
		materialDomain[i].Title = material.Title
		materialDomain[i].URL = material.URL
		materialDomain[i].Description = material.Description
		materialDomain[i].CreatedAt = material.CreatedAt
		materialDomain[i].UpdatedAt = material.UpdatedAt
	}

	progresses, _ := dc.menteeProgressRepository.FindByMentee(menteeId, courseId)

	for i := range materialDomain {
		for j := range progresses {
			if progresses[j].Completed && materialDomain[i].MaterialId == progresses[j].MaterialId {
				materialDomain[i].Completed = true
			}
		}
	}

	moduleDomain := make([]domain.DetailModule, len(modules))

	for i, module := range modules {
		moduleDomain[i].ModuleId = module.ID
		moduleDomain[i].CourseId = module.CourseId
		moduleDomain[i].Title = module.Title
		moduleDomain[i].Description = module.Description
		moduleDomain[i].CreatedAt = module.CreatedAt
		moduleDomain[i].UpdatedAt = module.UpdatedAt
	}

	for i, module := range moduleDomain {
		for _, material := range materialDomain {
			if module.ModuleId == material.ModuleId {
				moduleDomain[i].Materials = append(moduleDomain[i].Materials, material)
			}
		}
	}

	totalMaterialsArray, _ := dc.materialRepository.CountByCourse([]string{courseId})

	var totalMaterials int64

	if len(totalMaterialsArray) != 0 {
		totalMaterials = totalMaterialsArray[0]
	}

	if assignment != nil {
		totalMaterials += 1
	}

	progressArray, _ := dc.menteeProgressRepository.Count(menteeId, course.Title, menteeCourse.Status)

	var progress int64

	if len(progressArray) != 0 {
		progress = progressArray[0]
	}

	if isCompletingAssignment {
		progress += 1
	}

	courseDomain := domain.DetailCourse{
		CourseId:       course.ID,
		CategoryId:     course.CategoryId,
		MentorId:       course.MentorId,
		Mentor:         course.Mentor.Fullname,
		Category:       course.Category.Name,
		Title:          course.Title,
		Description:    course.Description,
		Thumbnail:      course.Thumbnail,
		TotalReviews:   course.TotalReviews,
		Rating:         course.Rating,
		Progress:       progress,
		TotalMaterials: totalMaterials,
		Modules:        moduleDomain,
		Assignment:     assignmentDomain,
		CreatedAt:      course.CreatedAt,
		UpdatedAt:      course.UpdatedAt,
	}

	return &courseDomain, nil
}
