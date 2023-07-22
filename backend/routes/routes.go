package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/superosystem/TrainingSystem/backend/helper"
	"github.com/superosystem/TrainingSystem/backend/middleware"

	_driverFactory "github.com/superosystem/TrainingSystem/backend/infrastructure"

	_menteeController "github.com/superosystem/TrainingSystem/backend/controller/mentees"
	_menteeUsecase "github.com/superosystem/TrainingSystem/backend/domain/mentees"

	_mentorController "github.com/superosystem/TrainingSystem/backend/controller/mentors"
	_mentorUsecase "github.com/superosystem/TrainingSystem/backend/domain/mentors"

	_otpController "github.com/superosystem/TrainingSystem/backend/controller/otp"
	_otpUsecase "github.com/superosystem/TrainingSystem/backend/domain/otp"

	_categoryController "github.com/superosystem/TrainingSystem/backend/controller/categories"
	_categoryUsecase "github.com/superosystem/TrainingSystem/backend/domain/categories"

	_courseController "github.com/superosystem/TrainingSystem/backend/controller/courses"
	_courseUsecase "github.com/superosystem/TrainingSystem/backend/domain/courses"

	_moduleController "github.com/superosystem/TrainingSystem/backend/controller/modules"
	_moduleUsecase "github.com/superosystem/TrainingSystem/backend/domain/modules"

	_assignmentController "github.com/superosystem/TrainingSystem/backend/controller/assignments"
	_assignmentUsecase "github.com/superosystem/TrainingSystem/backend/domain/assignments"

	_materialController "github.com/superosystem/TrainingSystem/backend/controller/materials"
	_materialUsecase "github.com/superosystem/TrainingSystem/backend/domain/materials"

	_menteeCoursesController "github.com/superosystem/TrainingSystem/backend/controller/menteeCourses"
	_menteeCoursesUsecase "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses"

	_menteeProgressController "github.com/superosystem/TrainingSystem/backend/controller/menteeProgresses"
	_menteeProgressesUsecase "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"

	_detailCourseController "github.com/superosystem/TrainingSystem/backend/controller/detailCourse"
	_detailCourseUsecase "github.com/superosystem/TrainingSystem/backend/domain/detailCourse"

	_assignmentMenteeController "github.com/superosystem/TrainingSystem/backend/controller/menteeAssignments"
	_assignmentMenteeUsecase "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"

	_manageMenteesController "github.com/superosystem/TrainingSystem/backend/controller/manageMentees"
	_manageMenteesUsecase "github.com/superosystem/TrainingSystem/backend/domain/manageMentees"

	_reviewController "github.com/superosystem/TrainingSystem/backend/controller/reviews"
	_reviewUsecase "github.com/superosystem/TrainingSystem/backend/domain/reviews"

	_certificateController "github.com/superosystem/TrainingSystem/backend/controller/certificates"
	_certificateUsecase "github.com/superosystem/TrainingSystem/backend/domain/certificates"
)

type RouteConfig struct {
	Echo          *echo.Echo            // echo top level instance
	MySQLDB       *gorm.DB              // mysql conn
	RedisDB       *redis.Client         // redis conn
	JWTConfig     *helper.JWTConfig     // JWT config
	Mailer        *helper.MailerConfig  // Mail config
	StorageConfig *helper.StorageConfig // Cloud storage config
}

func (routeConfig *RouteConfig) New() {
	// setup api v1
	v1 := routeConfig.Echo.Group("/api/v1")

	// setup auth middleware
	authMiddleware := middleware.NewAuthMiddleware(routeConfig.JWTConfig)

	// Inject the dependency to user
	userRepository := _driverFactory.NewUserRepository(routeConfig.MySQLDB)

	// Inject the dependency to otp
	otpRepository := _driverFactory.NewOTPRepository(routeConfig.RedisDB)
	otpUsecase := _otpUsecase.NewOTPUsecase(otpRepository, userRepository, routeConfig.Mailer)
	otpController := _otpController.NewOTPController(otpUsecase)

	// Inject the dependency to mentee
	menteeRepository := _driverFactory.NewMenteeRepository(routeConfig.MySQLDB)
	menteeUsecase := _menteeUsecase.NewMenteeUsecase(menteeRepository, userRepository, otpRepository, routeConfig.JWTConfig, routeConfig.Mailer, routeConfig.StorageConfig)
	menteeController := _menteeController.NewMenteeController(menteeUsecase, routeConfig.JWTConfig)

	// Inject the dependency to mentor
	mentorRepository := _driverFactory.NewMentorRepository(routeConfig.MySQLDB)
	mentorUsecase := _mentorUsecase.NewMentorUsecase(mentorRepository, userRepository, routeConfig.JWTConfig, routeConfig.StorageConfig, routeConfig.Mailer)
	mentorController := _mentorController.NewMentorController(mentorUsecase, routeConfig.JWTConfig)

	// Inject the dependency to category
	categoryRepository := _driverFactory.NewCategoryRepository(routeConfig.MySQLDB)
	categoryUsecase := _categoryUsecase.NewCategoryUsecase(categoryRepository)
	categoryController := _categoryController.NewCategoryController(categoryUsecase)

	// Inject the dependency to course
	courseRepository := _driverFactory.NewCourseRepository(routeConfig.MySQLDB)
	courseUsecase := _courseUsecase.NewCourseUsecase(courseRepository, mentorRepository, categoryRepository, routeConfig.StorageConfig)
	courseController := _courseController.NewCourseController(courseUsecase)

	// Inject the dependency to module
	moduleRepository := _driverFactory.NewModuleRepository(routeConfig.MySQLDB)
	moduleUsecase := _moduleUsecase.NewModuleUsecase(moduleRepository, courseRepository)
	moduleController := _moduleController.NewModuleController(moduleUsecase)

	// Inject the dependency to assignment
	assignmentRepository := _driverFactory.NewAssignmentRepository(routeConfig.MySQLDB)
	assignmentUsecase := _assignmentUsecase.NewAssignmentUsecase(assignmentRepository, courseRepository)
	assignmentController := _assignmentController.NewAssignmentsController(assignmentUsecase)

	// Inject the dependency to material
	materialRepository := _driverFactory.NewMaterialRepository(routeConfig.MySQLDB)
	materialUsecase := _materialUsecase.NewMaterialUsecase(materialRepository, moduleRepository, routeConfig.StorageConfig)
	materialController := _materialController.NewMaterialController(materialUsecase)

	// Inject the dependency to menteeProgress
	menteeProgressRepository := _driverFactory.NewMenteeProgressRepository(routeConfig.MySQLDB)
	menteeProgressUsecase := _menteeProgressesUsecase.NewMenteeProgressUsecase(menteeProgressRepository, menteeRepository, courseRepository, materialRepository)
	menteeProgressController := _menteeProgressController.NewMenteeProgressController(menteeProgressUsecase)

	// Inject the dependency to mentee assignment
	menteeAssignmentRepository := _driverFactory.NewMenteeAssignmentRepository(routeConfig.MySQLDB)
	menteeAssignmentUsecase := _assignmentMenteeUsecase.NewMenteeAssignmentUsecase(menteeAssignmentRepository, assignmentRepository, menteeRepository, routeConfig.StorageConfig)
	menteeAssignmentController := _assignmentMenteeController.NewAssignmentsMenteeController(menteeAssignmentUsecase, routeConfig.JWTConfig)

	// Inject the dependency to menteeCourse
	menteeCourseRepository := _driverFactory.NewMenteeCourseRepository(routeConfig.MySQLDB)
	menteeCourseUsecase := _menteeCoursesUsecase.NewMenteeCourseUsecase(menteeCourseRepository, menteeRepository, courseRepository, materialRepository, menteeProgressRepository, assignmentRepository, menteeAssignmentRepository)
	menteeCourseController := _menteeCoursesController.NewMenteeCourseController(menteeCourseUsecase)

	// Inject the dependency to detailCourse
	detailCourseUsecase := _detailCourseUsecase.NewDetailCourseUsecase(menteeRepository, courseRepository, moduleRepository, materialRepository, menteeProgressRepository, assignmentRepository, menteeAssignmentRepository, menteeCourseRepository)
	detailCourseController := _detailCourseController.NewDetailCourseController(detailCourseUsecase)

	// Inject the dependency to manageMentees
	manageMenteeUsecase := _manageMenteesUsecase.NewManageMenteeUsecase(menteeCourseRepository, menteeProgressRepository, menteeAssignmentRepository, routeConfig.StorageConfig)
	manageMenteeController := _manageMenteesController.NewManageMenteeController(manageMenteeUsecase)

	// Inject the dependency to review
	reviewRepository := _driverFactory.NewReviewRepository(routeConfig.MySQLDB)
	reviewUsecase := _reviewUsecase.NewReviewUsecase(reviewRepository, menteeCourseRepository, menteeRepository, courseRepository)
	reviewController := _reviewController.NewReviewController(reviewUsecase)

	// Inject the dependency to certificate
	certificateUsecase := _certificateUsecase.NewCertificateUsecase(menteeRepository, courseRepository)
	certificateController := _certificateController.NewCertificateController(certificateUsecase)

	// authentication routes
	auth := v1.Group("/auth")
	auth.POST("/mentee/login", menteeController.HandlerLoginMentee)
	auth.POST("/mentee/register", menteeController.HandlerRegisterMentee)
	auth.POST("/mentee/register/verify", menteeController.HandlerVerifyRegisterMentee)
	auth.POST("/mentee/forgot-password", menteeController.HandlerForgotPassword)
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
