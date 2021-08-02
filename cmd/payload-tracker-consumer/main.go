package main

import (
	"context"
	"log"

	"github.com/redhatinsights/payload-tracker-go/internal/config"
	"github.com/redhatinsights/payload-tracker-go/internal/kafka"
)

func main() {
	cfg := config.Get()
	ctx := context.Background()

	log.Println("Starting a new kafka consumer...")
	log.Println("Config for Consumer: ", cfg)

	consumer, err := kafka.NewConsumer(ctx, cfg, cfg.KafkaConfig.KafkaTopic)

	if err != nil {
		log.Println("ERROR! ", err)
	}

	// TODO: Add Handler in here

	kafka.NewConsumerEventLoop(ctx, consumer) // TODO: add in handler
}
