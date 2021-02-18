package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FileName   = "config"
	FileFormat = "yaml"
)

type configuration struct {
	Slack struct {
		WebhookUrl string `mapstructure:"webhook-url"`
	}
	SQS struct {
		DLQName string `mapstructure:"dlq-name"`
	}
	Internal struct {
		WorkerPool     int  `mapstructure:"worker-pool"`
		StructuredLogs bool `mapstructure:"structured-logs"`
	}
}

func LoadConfiguration(log *logrus.Logger) *configuration {
	pflag.String("slack.webhook-url", "", "slack webhook url")
	pflag.String("sqs.dlq-name", "", "sqs dead-letter queue name")
	pflag.Int("internal.worker-pool", 1, "the size of the worker pool")
	pflag.Bool("internal.structured-logs", false, "print logs using json format")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalf("Fatal error parsing flags: %v", err)
	}

	viper.SetConfigName(FileName)
	viper.SetConfigType(FileFormat)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Fatal error reading configuration file: %v", err)
		}
	}

	var config *configuration
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Fatal error decoding file: %v", err)
	}

	if config.Slack.WebhookUrl == "" {
		log.Fatal("Empty webhook URL")
	}

	if config.SQS.DLQName == "" {
		log.Fatal("Empty dead-letter queue name")
	}

	if config.Internal.WorkerPool < 1 {
		log.Fatal("Invalid worker-pool size")
	}

	return config
}
