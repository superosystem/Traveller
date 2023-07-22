package drivers

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	otpDomain "github.com/superosystem/TrainingSystem/backend/domain/otp"
	otpDB "github.com/superosystem/TrainingSystem/backend/infrastructure/redis/otp"

	menteeDomain "github.com/superosystem/TrainingSystem/backend/domain/mentees"
	menteeDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/mentees"

	mentorsDomain "github.com/superosystem/TrainingSystem/backend/domain/mentors"
	mentorsDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/mentors"

	userDomain "github.com/superosystem/TrainingSystem/backend/domain/users"
	userDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/users"

	categoryDomain "github.com/superosystem/TrainingSystem/backend/domain/categories"
	categoryDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/categories"

	courseDomain "github.com/superosystem/TrainingSystem/backend/domain/courses"
	courseDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/courses"

	moduleDomain "github.com/superosystem/TrainingSystem/backend/domain/modules"
	moduleDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/modules"

	assignmentDomain "github.com/superosystem/TrainingSystem/backend/domain/assignments"
	assignmentDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/assignments"

	materialDomain "github.com/superosystem/TrainingSystem/backend/domain/materials"
	materialDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/materials"

	menteeCoursesDomain "github.com/superosystem/TrainingSystem/backend/domain/menteeCourses"
	menteeCoursesDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/menteeCourses"

	menteeProgressesDomain "github.com/superosystem/TrainingSystem/backend/domain/menteeProgresses"
	menteeProgressesDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/menteeProgresses"

	menteeAssignmentsDomain "github.com/superosystem/TrainingSystem/backend/domain/menteeAssignments"
	menteeAssignmentsDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/menteeAssignments"

	reviewDomain "github.com/superosystem/TrainingSystem/backend/domain/reviews"
	reviewDB "github.com/superosystem/TrainingSystem/backend/infrastructure/mysql/reviews"
)

func NewOTPRepository(client *redis.Client) otpDomain.Repository {
	return otpDB.NewRedisRepository(client)
}

func NewUserRepository(conn *gorm.DB) userDomain.Repository {
	return userDB.NewMySQLRepository(conn)
}

func NewMenteeRepository(conn *gorm.DB) menteeDomain.Repository {
	return menteeDB.NewSQLRepository(conn)
}

func NewMentorRepository(conn *gorm.DB) mentorsDomain.Repository {
	return mentorsDB.NewSQLRepository(conn)
}

func NewCategoryRepository(conn *gorm.DB) categoryDomain.Repository {
	return categoryDB.NewSQLRepository(conn)
}

func NewCourseRepository(conn *gorm.DB) courseDomain.Repository {
	return courseDB.NewSQLRepository(conn)
}

func NewModuleRepository(conn *gorm.DB) moduleDomain.Repository {
	return moduleDB.NewSQLRepository(conn)
}

func NewAssignmentRepository(conn *gorm.DB) assignmentDomain.Repository {
	return assignmentDB.NewSQLRepository(conn)
}

func NewMaterialRepository(conn *gorm.DB) materialDomain.Repository {
	return materialDB.NewSQLRepository(conn)
}

func NewMenteeCourseRepository(conn *gorm.DB) menteeCoursesDomain.Repository {
	return menteeCoursesDB.NewSQLRepository(conn)
}

func NewMenteeProgressRepository(conn *gorm.DB) menteeProgressesDomain.Repository {
	return menteeProgressesDB.NewSQLRepository(conn)
}

func NewMenteeAssignmentRepository(conn *gorm.DB) menteeAssignmentsDomain.Repository {
	return menteeAssignmentsDB.NewSQLRepository(conn)
}

func NewReviewRepository(conn *gorm.DB) reviewDomain.Repository {
	return reviewDB.NewSQLRepository(conn)
}
