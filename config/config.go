package config

import (
  "github.com/sirupsen/logrus"
  "github.com/spf13/viper"
)

type configuration struct {
  Slack struct {
    WebhookUrl string `mapstructure:"webhook-url"`
  }
  SQS struct {
    DLQName string `mapstructure:"dlq-name"`
  }
  X9 struct {
    WorkerPool int `mapstructure:"worker-pool"`
  }
}

func LoadConfiguration(log *logrus.Logger) *configuration {
  viper.SetConfigName("config")
  viper.SetConfigType("yaml")
  viper.AddConfigPath(".")

  if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
      log.Fatalf("Configuration file not found: %v\n", err)
    } else {
      log.Fatalf("Fatal error reading configuration file: %v\n", err)
    }
  }

  viper.SetDefault("x9.worker-pool", 10)

  var config *configuration
  err := viper.Unmarshal(&config)
  if err != nil {
    log.Fatalf("Fatal error decoding file: %v\n", err)
  }

  if config.Slack.WebhookUrl == "" {
    log.Fatal("Empty webhook URL")
  }

  if config.SQS.DLQName == "" {
    log.Fatal("Empty dead-letter queue name")
  }

  return config
}
