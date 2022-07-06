package database

import (
	"final-project/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	err      error
	username string
	port     string
	password string
	dbName   string
	dsn      string
)

func StartDB() {
	var (
		username = os.Getenv("PG_USERNAME")
		port     = os.Getenv("PG_PORT")
		password = os.Getenv("PG_PASSWORD")
		dbName   = os.Getenv("PG_DATABASE_NAME")
	)

	dsn = fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", username, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	fmt.Println("Connection to the database is established")

	db.Debug().AutoMigrate(models.User{}, models.SocialMedia{}, models.Comment{}, models.Photo{})
}

func GetDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	return db
}
