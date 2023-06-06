package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/middlewares"
	"gorm.io/gorm"

	_repository "github.com/superosystem/trainingsystem-backend/src/repository"
	_use_case "github.com/superosystem/trainingsystem-backend/src/usecase"

	_assignmentController "github.com/superosystem/trainingsystem-backend/src/controllers/assignments"
	_categoryController "github.com/superosystem/trainingsystem-backend/src/controllers/categories"
	_certificateController "github.com/superosystem/trainingsystem-backend/src/controllers/certificates"
	_courseController "github.com/superosystem/trainingsystem-backend/src/controllers/courses"
	_detailCourseController "github.com/superosystem/trainingsystem-backend/src/controllers/detailCourse"
	_manageMenteesController "github.com/superosystem/trainingsystem-backend/src/controllers/manageMentees"
	_materialController "github.com/superosystem/trainingsystem-backend/src/controllers/materials"
	_assignmentMenteeController "github.com/superosystem/trainingsystem-backend/src/controllers/menteeAssignments"
	_menteeCoursesController "github.com/superosystem/trainingsystem-backend/src/controllers/menteeCourses"
	_menteeProgressController "github.com/superosystem/trainingsystem-backend/src/controllers/menteeProgresses"
	_menteeController "github.com/superosystem/trainingsystem-backend/src/controllers/mentees"
	_mentorController "github.com/superosystem/trainingsystem-backend/src/controllers/mentors"
	_moduleController "github.com/superosystem/trainingsystem-backend/src/controllers/modules"
	_otpController "github.com/superosystem/trainingsystem-backend/src/controllers/otp"
	_reviewController "github.com/superosystem/trainingsystem-backend/src/controllers/reviews"
)

type RouteConfig struct {
	// echo top level instance
	Echo *echo.Echo

	// mysql conn
	MySQLDB *gorm.DB

	// redis conn
	RedisDB *redis.Client

	// JWT config
	JWTConfig *config.JWTConfig

	// mail config
	Mailer *config.MailerConfig

	// cloud storage config
	StorageConfig *config.StorageConfig
}

func (routeConfig *RouteConfig) New() {
	// setup api v1
	v1 := routeConfig.Echo.Group("/api/v1")

	// setup auth middlewares
	authMiddleware := middlewares.NewAuthMiddleware(routeConfig.JWTConfig)

	// Inject the dependency to user
	userRepository := _repository.NewUserRepository(routeConfig.MySQLDB)

	// Inject the dependency to otp
	otpRepository := _repository.NewOtpRepository(routeConfig.RedisDB)
	otpUseCase := _use_case.NewOTPUseCase(otpRepository, userRepository, routeConfig.Mailer)
	otpController := _otpController.NewOTPController(otpUseCase)

	// Inject the dependency to mentee
	menteeRepository := _repository.NewMenteeRepository(routeConfig.MySQLDB)
	menteeUseCase := _use_case.NewMenteeUseCase(menteeRepository, userRepository, otpRepository, routeConfig.JWTConfig, routeConfig.Mailer, routeConfig.StorageConfig)
	menteeController := _menteeController.NewMenteeController(menteeUseCase, routeConfig.JWTConfig)

	// Inject the dependency to mentor
	mentorRepository := _repository.NewMentorRepository(routeConfig.MySQLDB)
	mentorUseCase := _use_case.NewMentorUseCase(mentorRepository, userRepository, routeConfig.JWTConfig, routeConfig.StorageConfig, routeConfig.Mailer)
	mentorController := _mentorController.NewMentorController(mentorUseCase, routeConfig.JWTConfig)

	// Inject the dependency to category
	categoryRepository := _repository.NewCategoryRepository(routeConfig.MySQLDB)
	categoryUseCase := _use_case.NewCategoryUseCase(categoryRepository)
	categoryController := _categoryController.NewCategoryController(categoryUseCase)

	// Inject the dependency to course
	courseRepository := _repository.NewCourseRepository(routeConfig.MySQLDB)
	courseUseCase := _use_case.NewCourseUseCase(courseRepository, mentorRepository, categoryRepository, routeConfig.StorageConfig)
	courseController := _courseController.NewCourseController(courseUseCase)

	// Inject the dependency to module
	moduleRepository := _repository.NewModuleRepository(routeConfig.MySQLDB)
	moduleUseCase := _use_case.NewModuleUseCase(moduleRepository, courseRepository)
	moduleController := _moduleController.NewModuleController(moduleUseCase)

	// Inject the dependency to assignment
	assignmentRepository := _repository.NewAssignmentRepository(routeConfig.MySQLDB)
	assignmentUseCase := _use_case.NewAssignmentUseCase(assignmentRepository, courseRepository)
	assignmentController := _assignmentController.NewAssignmentsController(assignmentUseCase)

	// Inject the dependency to material
	materialRepository := _repository.NewMaterialRepository(routeConfig.MySQLDB)
	materialUseCase := _use_case.NewMaterialUseCase(materialRepository, moduleRepository, routeConfig.StorageConfig)
	materialController := _materialController.NewMaterialController(materialUseCase)

	// Inject the dependency to menteeProgress
	menteeProgressRepository := _repository.NewMenteeProgressRepository(routeConfig.MySQLDB)
	menteeProgressUseCase := _use_case.NewMenteeProgressUseCase(menteeProgressRepository, menteeRepository, courseRepository, materialRepository)
	menteeProgressController := _menteeProgressController.NewMenteeProgressController(menteeProgressUseCase)

	// Inject the dependency to mentee assignment
	menteeAssignmentRepository := _repository.NewMenteeAssignmentRepository(routeConfig.MySQLDB)
	menteeAssignmentUseCase := _use_case.NewMenteeAssignmentUseCase(menteeAssignmentRepository, assignmentRepository, menteeRepository, routeConfig.StorageConfig)
	menteeAssignmentController := _assignmentMenteeController.NewAssignmentsMenteeController(menteeAssignmentUseCase, routeConfig.JWTConfig)

	// Inject the dependency to menteeCourse
	menteeCourseRepository := _repository.NewMenteeCourseRepository(routeConfig.MySQLDB)
	menteeCourseUseCase := _use_case.NewMenteeCourseUseCase(menteeCourseRepository, menteeRepository, courseRepository, materialRepository, menteeProgressRepository, assignmentRepository, menteeAssignmentRepository)
	menteeCourseController := _menteeCoursesController.NewMenteeCourseController(menteeCourseUseCase)

	detailCourseUseCase := _use_case.NewDetailCourseUseCase(menteeRepository, courseRepository, moduleRepository, materialRepository, menteeProgressRepository, assignmentRepository, menteeAssignmentRepository, menteeCourseRepository)
	detailCourseController := _detailCourseController.NewDetailCourseController(detailCourseUseCase)

	manageMenteeUseCase := _use_case.NewManageMenteeUseCase(menteeCourseRepository, menteeProgressRepository, menteeAssignmentRepository, routeConfig.StorageConfig)
	manageMenteeController := _manageMenteesController.NewManageMenteeController(manageMenteeUseCase)

	reviewRepository := _repository.NewReviewRepository(routeConfig.MySQLDB)
	reviewUseCase := _use_case.NewReviewUseCase(reviewRepository, menteeCourseRepository, menteeRepository, courseRepository)
	reviewController := _reviewController.NewReviewController(reviewUseCase)

	certificateUseCase := _use_case.NewCertificateUseCase(menteeRepository, courseRepository)
	certificateController := _certificateController.NewCertificateController(certificateUseCase)

	// authentication routes
	auth := v1.Group("/auth")
	auth.POST("/mentee/login", menteeController.HandlerLoginMentee)
	auth.POST("/mentee/register", menteeController.HandlerRegisterMentee)
	auth.POST("/mentee/register/verify", menteeController.HandlerVerifyRegisterMentee)
	auth.POST("/forgot-password", menteeController.HandlerForgotPassword)
	auth.POST("/send-otp", otpController.HandlerSendOTP)
	auth.POST("/check-otp", otpController.HandlerCheckOTP)
	auth.POST("/mentor/login", mentorController.HandlerLoginMentor)
	auth.POST("/mentor/register", mentorController.HandlerRegisterMentor)
	auth.POST("/mentor/forgot-password", mentorController.HandlerForgotPassword)

	// mentor routes
	mentor := v1.Group("/mentors", authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	mentor.GET("", mentorController.HandlerFindAll)
	mentor.GET("/profile", mentorController.HandlerProfileMentor)
	mentor.PUT("/:mentorId/update-password", mentorController.HandlerUpdatePassword)
	mentor.GET("/:mentorId", mentorController.HandlerFindByID)
	mentor.PUT("/:mentorId", mentorController.HandlerUpdateProfile)

	// mentee routes
	mentee := v1.Group("/mentees", authMiddleware.IsAuthenticated(), authMiddleware.IsMentee)
	mentee.GET("", menteeController.HandlerFindAll)
	mentee.POST("/progress", menteeProgressController.HandlerAddProgress)
	mentee.GET("/profile", menteeController.HandlerProfileMentee)
	mentee.GET("/:menteeId", menteeController.HandlerFindByID)
	mentee.PUT("/:menteeId", menteeController.HandlerUpdateProfile)
	mentee.GET("/:menteeId/reviews", reviewController.HandlerFindByMentee)
	mentee.GET("/:menteeId/courses", menteeCourseController.HandlerFindMenteeCourses)
	mentee.GET("/:menteeId/courses/:courseId/certificate", certificateController.HandlerGenerateCert)
	mentee.GET("/:menteeId/courses/:courseId/details", detailCourseController.HandlerDetailCourseEnrolled)
	mentee.PUT("/:menteeId/courses/:courseId/complete", menteeCourseController.HandlerCompleteCourse)
	mentee.GET("/:menteeId/courses/:courseId", menteeCourseController.HandlerCheckEnrollmentCourse)
	mentee.GET("/:menteeId/materials/:materialId", menteeProgressController.HandlerFindMaterialEnrolled)

	//	category routes
	cat := v1.Group("/categories")
	cat.POST("", categoryController.HandlerCreateCategory, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	cat.GET("", categoryController.HandlerFindAllCategories)
	cat.GET("/:categoryId", categoryController.HandlerFindByIdCategory)
	cat.PUT("/:categoryId", categoryController.HandlerUpdateCategory, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)

	// course routes
	course := v1.Group("/courses")
	course.POST("", courseController.HandlerCreateCourse, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	course.GET("", courseController.HandlerFindAllCourses)
	course.GET("/popular", courseController.HandlerFindByPopular)
	course.POST("/enroll-course", menteeCourseController.HandlerEnrollCourse, authMiddleware.IsAuthenticated())
	course.GET("/categories/:categoryId", courseController.HandlerFindByCategory)
	course.GET("/mentors/:mentorId", courseController.HandlerFindByMentor)
	course.GET("/:courseId/reviews", reviewController.HandlerFindByCourse)
	course.GET("/:courseId/mentees", menteeController.HandlerFindMenteesByCourse)
	course.DELETE("/:courseId/mentees/:menteeId/delete-access", manageMenteeController.HandlerDeleteAccessMentee)
	course.GET("/:courseId/details", detailCourseController.HandlerDetailCourse)
	course.GET("/:courseId", courseController.HandlerFindByIdCourse)
	course.PUT("/:courseId", courseController.HandlerUpdateCourse, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	course.DELETE("/:courseId", courseController.HandlerSoftDeleteCourse, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)

	// module routes
	module := v1.Group("/modules")
	module.POST("", moduleController.HandlerCreateModule, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	module.GET("/:moduleId", moduleController.HandlerFindByIdModule)
	module.PUT("/:moduleId", moduleController.HandlerUpdateModule)
	module.DELETE("/:moduleId", moduleController.HandlerDeleteModule)
	module.PUT("/:moduleId", moduleController.HandlerUpdateModule, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	module.DELETE("/:moduleId", moduleController.HandlerDeleteModule, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)

	// assignment routes
	assignment := v1.Group("/assignments")
	assignment.POST("", assignmentController.HandlerCreateAssignment, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	assignment.GET("/:assignmentId", assignmentController.HandlerFindByIdAssignment)
	assignment.GET("/courses/:courseid", assignmentController.HandlerFindByCourse)
	assignment.PUT("/:assignmentId", assignmentController.HandlerUpdateAssignment, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	assignment.DELETE("/:assignmentId", assignmentController.HandlerDeleteAssignment, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)

	// material routes
	material := v1.Group("/materials")
	material.POST("", materialController.HandlerCreateMaterial, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	material.DELETE("/modules/:moduleId", materialController.HandlerSoftDeleteMaterialByModule, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	material.GET("/:materialId", materialController.HandlerFindByIdMaterial)
	material.PUT("/:materialId", materialController.HandlerUpdateMaterial, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	material.DELETE("/:materialId", materialController.HandlerSoftDeleteMaterial, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)

	// Mentee assignment routes
	menteeAssignment := v1.Group("/mentee-assignments")
	menteeAssignment.POST("", menteeAssignmentController.HandlerCreateMenteeAssignment, authMiddleware.IsAuthenticated(), authMiddleware.IsMentee)
	menteeAssignment.PUT("/:menteeAssignmentId", menteeAssignmentController.HandlerUpdateMenteeAssignment, authMiddleware.IsAuthenticated(), authMiddleware.IsMentee)
	menteeAssignment.GET("/:menteeAssignmentId", menteeAssignmentController.HandlerFindByIdMenteeAssignment, authMiddleware.IsAuthenticated())
	menteeAssignment.PUT("/grade/:menteeAssignmentId", menteeAssignmentController.HandlerUpdateGradeMentee, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	menteeAssignment.DELETE("/:menteeAssignmentId", menteeAssignmentController.HandlerSoftDeleteMenteeAssignment, authMiddleware.IsAuthenticated(), authMiddleware.IsMentee)
	menteeAssignment.GET("/assignments/:assignmentId", menteeAssignmentController.HandlerFindByAssignmentId, authMiddleware.IsAuthenticated(), authMiddleware.IsMentor)
	menteeAssignment.GET("/mentee", menteeAssignmentController.HandlerFindByMenteeId, authMiddleware.IsAuthenticated(), authMiddleware.IsMentee)
	menteeAssignment.GET("/:menteeId/assignments/:assignmentId", menteeAssignmentController.HandlerFindMenteeAssignmentEnrolled, authMiddleware.IsAuthenticated())

	// reviews routes
	review := v1.Group("/reviews", authMiddleware.IsAuthenticated(), authMiddleware.IsMentee)
	review.POST("", reviewController.HandlerCreateReview)
}
