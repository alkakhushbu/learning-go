package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Video struct {
	      gorm.Model
	Title       string
	Description string
	Loc         string
}

func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error opening DB connection", err)
		return
	}
	// create a table according the go struct
	// err = db.AutoMigrate(&Video{})
	// if err != nil {
	// 	panic("failed to migrate table")
	// }
	log.Println("Table Video created")
	video := Video{Title: "Learning Go",
		Description: "An Idiomatic Approach to Real World Go Programming",
		Loc:         "AWS S3"}

	err = db.Create(&video).Error
	if err != nil {
		log.Println("Error in creating entry in DB:", err)
		return
	}
	err = db.First(&video, 1).Error
	if err != nil {
		log.Println("Cannot find video :", err)
		return
	}
}
