package config

import (
	"log"
	"os"

	"example.com/event-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("Environment variable belum diisi")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("gagal terkoneksi di database", err)
	}

	err = database.AutoMigrate(&models.Event{})
	if err != nil {
		log.Fatal("Gagal melakukan migration database : ", err)
	}

	DB = database
	log.Println("Berhasil terkoneksi ke Database")
}
