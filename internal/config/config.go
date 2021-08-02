package config

import (
	"strings"

	"github.com/spf13/viper"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

var rdsCaPath *string

// Get sets each config option with its defaults
func Get() *viper.Viper {
	options := viper.New()

	options.SetDefault("db.user", "crc")
	options.SetDefault("db.password", "crc")
	options.SetDefault("db.name", "crc")
	options.SetDefault("db.host", "0.0.0.0")
	options.SetDefault("db.port", "5432")

	options.SetDefault("kafka.timeout", 10000)
	options.SetDefault("kafka.group.id", "payload-tracker-go")
	options.SetDefault("kafka.auto.offset.reset", "latest")
	options.SetDefault("kafka.auto.commit.interval.ms", 5000)
	options.SetDefault("kafka.request.required.acks", -1) // -1 == "all"
	options.SetDefault("kafka.message.send.max.retries", 15)
	options.SetDefault("kafka.retry.backoff.ms", 100)

	if clowder.IsClowderEnabled() {
		options.SetDefault("kafka.brokers", clowder.KafkaServers)
		options.SetDefault("topic.payload.status", clowder.KafkaTopics["platform.payload-status"].Name)
	} else {
		options.SetDefault("kafka.brokers", []string{"localhost:29092"})
		options.SetDefault("topic.payload.status", "platform.payload-status")
	}

	options.AutomaticEnv()
	options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return options
}
