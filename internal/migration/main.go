package main

import (
	"github.com/redhatinsights/payload-tracker-go/internal/db"
	"github.com/redhatinsights/payload-tracker-go/internal/logging"
	"github.com/redhatinsights/payload-tracker-go/internal/models"
)

func main() {
	logging.InitLogger()

	db.DbConnect()

	db.DB.AutoMigrate(
		&models.Services{},
		&models.Sources{},
		&models.Statuses{},
		&models.PayloadStatuses{},
		&models.Payloads{},
	)

	logging.Log.Info("DB Migration Complete")
}
