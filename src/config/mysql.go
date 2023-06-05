package config

import (
	"fmt"
	"log"

	"github.com/superosystem/trainingsystem-backend/src/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigMySQL struct {
	MYSQL_USERNAME string
	MYSQL_PASSWORD string
	MYSQL_NAME     string
	MYSQL_HOST     string
	MYSQL_PORT     string
}

func (config *ConfigMySQL) InitMySQLDatabase() *gorm.DB {
	var err error

	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MYSQL_USERNAME,
		config.MYSQL_PASSWORD,
		config.MYSQL_HOST,
		config.MYSQL_PORT,
		config.MYSQL_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("error when connecting to the database: %s", err)
	}

	log.Println("connected to the database")

	return db
}

func MySQLMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(
		&entities.User{},
		&entities.Mentee{},
		&entities.Mentor{},
		&entities.Category{},
		&entities.Course{},
		&entities.Module{},
		&entities.Material{},
		&entities.Assignment{},
		&entities.MenteeCourse{},
		&entities.MenteeProgress{},
		&entities.MenteeAssignment{},
		&entities.Review{},
	)
}

func CloseDB(db *gorm.DB) error {
	database, err := db.DB()

	if err != nil {
		log.Printf("error when getting the database instance: %v", err)
		return err
	}

	if err := database.Close(); err != nil {
		log.Printf("error when closing the database connection: %v", err)
		return err
	}

	log.Println("database connection is closed")

	return nil
}
