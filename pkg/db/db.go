package db

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/redhatinsights/payload-tracker-go/internal/config"
	"github.com/redhatinsights/payload-tracker-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var cfg *viper.Viper = config.Get()

var (
	user = cfg.GetString("db.user")
	password = cfg.GetString("db.password")
	dbname = cfg.GetString("db.name")
	host = cfg.GetString("db.host")
	port = cfg.GetString("db.port")
)

func DbConnect() {


	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.PayloadStatuses{},
		&models.Payloads{},
		&models.Services{},
		&models.Sources{},
		&models.Statuses{},
	)

	DB = db

	fmt.Println("DB initialization complete")
}
