package models

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	configurations := "parseTime=true"
	dbUrlConnnect := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?" + configurations
	db, err := gorm.Open(mysql.Open(dbUrlConnnect))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Student{}, &Instructor{}, &Category{}, &Learning_path{}, &Expertise{}, &Lesson{}, &Sayembara{}, &Rating{})
	DB = db
}
