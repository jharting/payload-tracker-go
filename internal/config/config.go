package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

type TrackerConfig struct {
	PublicPort       string
	MetricsPort      string
	LogLevel string
	Hostname string
	KafkaConfig      KafkaCfg
	CloudwatchConfig CloudwatchCfg
	DatabaseConfig   DatabaseCfg
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

type DatabaseCfg struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	RDSCa      string
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

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// global logging
	options.SetDefault("logLevel", "INFO")
	options.SetDefault("Hostname", hostname)

	// kafka config
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
		// database
		options.SetDefault("db.user", cfg.Database.Username)
		options.SetDefault("db.password", cfg.Database.Password)
		options.SetDefault("db.name", cfg.Database.Name)
		options.SetDefault("db.host", cfg.Database.Hostname)
		options.SetDefault("db.port", cfg.Database.Port)
		options.SetDefault("rdsCa", cfg.Database.RdsCa)
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
		// database
		options.SetDefault("db.user", "crc")
		options.SetDefault("db.password", "crc")
		options.SetDefault("db.name", "crc")
		options.SetDefault("db.host", "0.0.0.0")
		options.SetDefault("db.port", "5432")
		// cloudwatch
		options.SetDefault("logGroup", "platform-dev")
		options.SetDefault("cwRegion", "us-east-1")
		options.SetDefault("cwAccessKey", os.Getenv("CW_AWS_ACCESS_KEY_ID"))
		options.SetDefault("cwSecretKey", os.Getenv("CW_AWS_SECRET_ACCESS_KEY"))
	}

	options.AutomaticEnv()
	options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	trackerCfg := &TrackerConfig{
		LogLevel: options.GetString("logLevel"),
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
		DatabaseConfig: DatabaseCfg{
			DBUser:     options.GetString("db.user"),
			DBPassword: options.GetString("db.password"),
			DBName:     options.GetString("db.name"),
			DBHost:     options.GetString("db.host"),
			DBPort:     options.GetString("db.port"),
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

			// write the Kafka CA path using the app-common-go package
			caPath, err := cfg.KafkaCa(broker)

			if err != nil {
				panic("Kafka CA Failed to Write")
			}

			trackerCfg.KafkaConfig.KafkaCA = caPath

		}

		// write the RDS CA using the app-common-go package
		rdsCAPath, err := cfg.RdsCa()

		if err != nil {
			panic("RDS CA Failed to Write")
		}

		trackerCfg.DatabaseConfig.RDSCa = rdsCAPath
	}

	return trackerCfg
}
