package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	databaseName := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v dbname=%v port=%v", host, user, databaseName, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Post{}, &User{})

	DB = database
}
