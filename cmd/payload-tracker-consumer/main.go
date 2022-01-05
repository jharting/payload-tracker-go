package main

import (
	"context"

	"github.com/redhatinsights/payload-tracker-go/internal/config"
	"github.com/redhatinsights/payload-tracker-go/internal/db"
	"github.com/redhatinsights/payload-tracker-go/internal/kafka"
	"github.com/redhatinsights/payload-tracker-go/internal/logging"
)

func main() {
	logging.InitLogger()

	cfg := config.Get()
	ctx := context.Background()

	logging.Log.Info("Setting up DB")
	db.DbConnect(cfg)

	logging.Log.Info("Starting a new kafka consumer...")

	consumer, err := kafka.NewConsumer(ctx, cfg, cfg.KafkaConfig.KafkaTopic)

	if err != nil {
		logging.Log.Fatal("ERROR! ", err)
	}

	kafka.NewConsumerEventLoop(ctx, cfg, consumer, db.DB)
}
