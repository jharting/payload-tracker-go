package config

import (
	"strings"

	"github.com/spf13/viper"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

var rdsCaPath *string

type TrackerConfig struct {
	PublicPort string
	MetricsPort string
	KafkaConfig KafkaCfg
}

type KafkaCfg struct {
	KafkaTimeout int
	KafkaGroupID string
	KafkaAutoOffsetReset string
	KafkaAutoCommitInterval int
	KafkaRequestRequiredAcks int
	KafkaMessageSendMaxRetries int
	KafkaRetryBackoffMs int
	KafkaBrokers []string
	KafkaTopic string
	KafkaUsername string
	KafkaPassword string
	KafkaCA string
	SASLMechanism string
	Protocol string
}

// Get sets each config option with its defaults
func Get() *TrackerConfig {
	options := viper.New()

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
		options.SetDefault("publicPort", clowder.LoadedConfig.PublicPort)
		options.SetDefault("metricsPort", clowder.LoadedConfig.MetricsPort)
	} else {
		options.SetDefault("kafka.brokers", []string{"localhost:29092"})
		options.SetDefault("topic.payload.status", "platform.payload-status")
		options.SetDefault("publicPort", "8080")
		options.SetDefault("metricsPort", "8081")
	}

	options.AutomaticEnv()
	options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	trackerCfg := &TrackerConfig{
		PublicPort: options.GetString("publicPort"),
		MetricsPort: options.GetString("metricsPort"),
		KafkaConfig: KafkaCfg{
			KafkaTimeout: options.GetInt("kafka.timeout"),
			KafkaGroupID: options.GetString("kafka.group.id"),
			KafkaAutoOffsetReset: options.GetString("kafka.auto.offset.reset"),
			KafkaAutoCommitInterval: options.GetInt("kafka.auto.commit.interval.ms"),
			KafkaRequestRequiredAcks: options.GetInt("kafka.request.required.acks"),
			KafkaMessageSendMaxRetries: options.GetInt("kafka.message.send.max.retries"),
			KafkaRetryBackoffMs: options.GetInt("kafka.retry.backoff.ms"),
			KafkaBrokers: options.GetStringSlice("kafka.brokers"),
			KafkaTopic: options.GetString("topic.payload.status"),
		},
	}

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig
		broker := cfg.Kafka.Brokers[0]

		if broker.Authtype != nil {
			trackerCfg.KafkaConfig.KafkaUsername = *broker.Sasl.Username
			trackerCfg.KafkaConfig.KafkaPassword = *broker.Sasl.Password
			trackerCfg.KafkaConfig.SASLMechanism = "SCRAM-SHA-512"
			trackerCfg.KafkaConfig.Protocol = "sasl_ssl"
			caPath, err := cfg.KafkaCa(broker)

			if err != nil {
				panic("Kafka CA Failed to Write")
			}

			trackerCfg.KafkaConfig.KafkaCA = caPath

		}
	}

	return trackerCfg
}