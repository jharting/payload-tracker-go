package db

import (
	"fmt"
	"log"

	"github.com/redhatinsights/payload-tracker-go/internal/config"
	"github.com/redhatinsights/payload-tracker-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var cfg *config.TrackerConfig = config.Get()

var (
	user     = cfg.DatabaseConfig.DBUser
	password = cfg.DatabaseConfig.DBPassword
	dbname   = cfg.DatabaseConfig.DBName
	host     = cfg.DatabaseConfig.DBHost
	port     = cfg.DatabaseConfig.DBPort
)

func DbConnect() {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.Services{},
		&models.Sources{},
		&models.Statuses{},
		&models.PayloadStatuses{},
		&models.Payloads{},
	)
	// db.Model(&models.PayloadStatuses{}).AddForeignKey("source_id", "sources(id)", "RESTRICT", "RESTRICT")

	DB = db

	log.Println("DB initialization complete")
}
