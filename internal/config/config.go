package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

var rdsCaPath *string

type TrackerConfig struct {
	PublicPort  string
	MetricsPort string
	KafkaConfig KafkaCfg
	CloudwatchConfig CloudwatchCfg
}

type KafkaCfg struct {
	KafkaTimeout               int
	KafkaGroupID               string
	KafkaAutoOffsetReset       string
	KafkaAutoCommitInterval    int
	KafkaRequestRequiredAcks   int
	KafkaMessageSendMaxRetries int
	KafkaRetryBackoffMs        int
	KafkaBrokers               []string
	KafkaTopic                 string
	KafkaUsername              string
	KafkaPassword              string
	KafkaCA                    string
	SASLMechanism              string
	Protocol                   string
}

type CloudwatchCfg struct {
	CWLogGroup  string
	CWRegion    string
	CWAccessKey string
	CWSecretKey string
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
		cfg := clowder.LoadedConfig

		// kafka
		options.SetDefault("kafka.brokers", clowder.KafkaServers)
		options.SetDefault("topic.payload.status", clowder.KafkaTopics["platform.payload-status"].Name)
		// ports
		options.SetDefault("publicPort", cfg.PublicPort)
		options.SetDefault("metricsPort", cfg.MetricsPort)
		// cloudwatch
		options.SetDefault("logGroup", cfg.Logging.Cloudwatch.LogGroup)
		options.SetDefault("cwRegion", cfg.Logging.Cloudwatch.Region)
		options.SetDefault("cwAccessKey", cfg.Logging.Cloudwatch.AccessKeyId)
		options.SetDefault("cwSecretKey", cfg.Logging.Cloudwatch.SecretAccessKey)

	} else {
		
		// kafka
		options.SetDefault("kafka.brokers", []string{"localhost:29092"})
		options.SetDefault("topic.payload.status", "platform.payload-status")
		// ports
		options.SetDefault("publicPort", "8080")
		options.SetDefault("metricsPort", "8081")
		// cloudwatch
		options.SetDefault("logGroup", "platform-dev")
		options.SetDefault("cwRegion", "us-east-1")
		options.SetDefault("cwAccessKey", os.Getenv("CW_AWS_ACCESS_KEY_ID"))
		options.SetDefault("cwSecretKey", os.Getenv("CW_AWS_SECRET_ACCESS_KEY"))
	}

	options.AutomaticEnv()
	options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	trackerCfg := &TrackerConfig{
		PublicPort:  options.GetString("publicPort"),
		MetricsPort: options.GetString("metricsPort"),
		KafkaConfig: KafkaCfg{
			KafkaTimeout:               options.GetInt("kafka.timeout"),
			KafkaGroupID:               options.GetString("kafka.group.id"),
			KafkaAutoOffsetReset:       options.GetString("kafka.auto.offset.reset"),
			KafkaAutoCommitInterval:    options.GetInt("kafka.auto.commit.interval.ms"),
			KafkaRequestRequiredAcks:   options.GetInt("kafka.request.required.acks"),
			KafkaMessageSendMaxRetries: options.GetInt("kafka.message.send.max.retries"),
			KafkaRetryBackoffMs:        options.GetInt("kafka.retry.backoff.ms"),
			KafkaBrokers:               options.GetStringSlice("kafka.brokers"),
			KafkaTopic:                 options.GetString("topic.payload.status"),
		},
		CloudwatchConfig: CloudwatchCfg{
			CWLogGroup:  options.GetString("logGroup"),
			CWRegion:    options.GetString("cwRegion"),
			CWAccessKey: options.GetString("cwAccessKey"),
			CWSecretKey: options.GetString("cwSecretKey"),
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
