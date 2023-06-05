package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/superosystem/trainingsystem-backend/src/config"

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

	_use_case "github.com/superosystem/trainingsystem-backend/src/usecase"

	_driver_mysql "github.com/superosystem/trainingsystem-backend/src/drivers/mysql"
	_driver_redis "github.com/superosystem/trainingsystem-backend/src/drivers/redis"

	"github.com/superosystem/trainingsystem-backend/src/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

	// setup auth middleware
	authMiddleware := middleware.NewAuthMiddleware(routeConfig.JWTConfig)

	// Inject the dependency to user
	userRepository := _driver_mysql.NewUserRepository(routeConfig.MySQLDB)

	// Inject the dependency to otp
	otpRepository := _driver_redis.NewOtpRepository(routeConfig.RedisDB)
	otpUsecase := _use_case.NewOTPUsecase(otpRepository, userRepository, routeConfig.Mailer)
	otpController := _otpController.NewOTPController(otpUsecase)

	// Inject the dependency to mentee
	menteeRepository := _driver_mysql.NewMenteeRepository(routeConfig.MySQLDB)
	menteeUsecase := _use_case.NewMenteeUsecase(menteeRepository, userRepository, otpRepository, routeConfig.JWTConfig, routeConfig.Mailer, routeConfig.StorageConfig)
	menteeController := _menteeController.NewMenteeController(menteeUsecase, routeConfig.JWTConfig)

	// Inject the dependency to mentor
	mentorRepository := _driver_mysql.NewMentorRepository(routeConfig.MySQLDB)
	mentorUsecase := _use_case.NewMentorUsecase(mentorRepository, userRepository, routeConfig.JWTConfig, routeConfig.StorageConfig, routeConfig.Mailer)
	mentorController := _mentorController.NewMentorController(mentorUsecase, routeConfig.JWTConfig)

	// Inject the dependency to category
	categoryRepository := _driver_mysql.NewCategoryRepository(routeConfig.MySQLDB)
	categoryUsecase := _use_case.NewCategoryUsecase(categoryRepository)
	categoryController := _categoryController.NewCategoryController(categoryUsecase)

	// Inject the dependency to course
	courseRepository := _driver_mysql.NewCourseRepository(routeConfig.MySQLDB)
	courseUsecase := _use_case.NewCourseUsecase(courseRepository, mentorRepository, categoryRepository, routeConfig.StorageConfig)
	courseController := _courseController.NewCourseController(courseUsecase)

	// Inject the dependency to module
	moduleRepository := _driver_mysql.NewModuleRepository(routeConfig.MySQLDB)
	moduleUsecase := _use_case.NewModuleUsecase(moduleRepository, courseRepository)
	moduleController := _moduleController.NewModuleController(moduleUsecase)

	// Inject the dependency to assignment
	assignmentRepository := _driver_mysql.NewAssignmentRepository(routeConfig.MySQLDB)
	assignmentUsecase := _use_case.NewAssignmentUsecase(assignmentRepository, courseRepository)
	assignmentController := _assignmentController.NewAssignmentsController(assignmentUsecase)

	// Inject the dependency to material
	materialRepository := _driver_mysql.NewMaterialRepository(routeConfig.MySQLDB)
	materialUsecase := _use_case.NewMaterialUsecase(materialRepository, moduleRepository, routeConfig.StorageConfig)
	materialController := _materialController.NewMaterialController(materialUsecase)

	// Inject the dependency to menteeProgress
	menteeProgressRepository := _driver_mysql.NewMenteeProgressRepository(routeConfig.MySQLDB)
	menteeProgressUsecase := _use_case.NewMenteeProgressUsecase(menteeProgressRepository, menteeRepository, courseRepository, materialRepository)
	menteeProgressController := _menteeProgressController.NewMenteeProgressController(menteeProgressUsecase)

	// Inject the dependency to mentee assignment
	menteeAssignmentRepository := _driver_mysql.NewMenteeAssignmentRepository(routeConfig.MySQLDB)
	menteeAssignmentUsecase := _use_case.NewMenteeAssignmentUsecase(menteeAssignmentRepository, assignmentRepository, menteeRepository, routeConfig.StorageConfig)
	menteeAssignmentController := _assignmentMenteeController.NewAssignmentsMenteeController(menteeAssignmentUsecase, routeConfig.JWTConfig)

	// Inject the dependency to menteeCourse
	menteeCourseRepository := _driver_mysql.NewMenteeCourseRepository(routeConfig.MySQLDB)
	menteeCourseUsecase := _use_case.NewMenteeCourseUsecase(menteeCourseRepository, menteeRepository, courseRepository, materialRepository, menteeProgressRepository, assignmentRepository, menteeAssignmentRepository)
	menteeCourseController := _menteeCoursesController.NewMenteeCourseController(menteeCourseUsecase)

	detailCourseUsecase := _use_case.NewDetailCourseUsecase(menteeRepository, courseRepository, moduleRepository, materialRepository, menteeProgressRepository, assignmentRepository, menteeAssignmentRepository, menteeCourseRepository)
	detailCourseController := _detailCourseController.NewDetailCourseController(detailCourseUsecase)

	manageMenteeUsecase := _use_case.NewManageMenteeUsecase(menteeCourseRepository, menteeProgressRepository, menteeAssignmentRepository, routeConfig.StorageConfig)
	manageMenteeController := _manageMenteesController.NewManageMenteeController(manageMenteeUsecase)

	reviewRepository := _driver_mysql.NewReviewRepository(routeConfig.MySQLDB)
	reviewUsecase := _use_case.NewReviewUsecase(reviewRepository, menteeCourseRepository, menteeRepository, courseRepository)
	reviewController := _reviewController.NewReviewController(reviewUsecase)

	certificateUsecase := _use_case.NewCertificateUsecase(menteeRepository, courseRepository)
	certificateController := _certificateController.NewCertificateController(certificateUsecase)

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
