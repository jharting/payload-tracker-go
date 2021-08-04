package main

import (
	"context"

	"github.com/redhatinsights/payload-tracker-go/internal/config"
	"github.com/redhatinsights/payload-tracker-go/internal/kafka"
	"github.com/redhatinsights/payload-tracker-go/internal/logging"
)

func main() {
	logging.InitLogger()

	cfg := config.Get()
	ctx := context.Background()

	logging.Log.Info("Starting a new kafka consumer...")
	logging.Log.Info("Config for Consumer: ", cfg)

	consumer, err := kafka.NewConsumer(ctx, cfg, cfg.KafkaConfig.KafkaTopic)

	if err != nil {
		logging.Log.Fatal("ERROR! ", err)
	}

	// TODO: Add Handler in here

	kafka.NewConsumerEventLoop(ctx, consumer) // TODO: add in handler
}
