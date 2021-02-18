package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hpedrorodrigues/dlq-x9/config"
	"github.com/hpedrorodrigues/dlq-x9/consumer"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var logger = logrus.New()

	conf := config.LoadConfiguration(logger)

	if conf.Internal.StructuredLogs {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	cons := consumer.New(conf.SQS.DLQName, conf.Internal.WorkerPool)

	logger.Info("DLQ-X9")

	cons.Consume(logger, func(message *sqs.Message, log *logrus.Entry) error {
		log.Infof("Sending message to Slack channel [Id: %s]", *message.MessageId)

		text := fmt.Sprintf(
			"Hey, a new message was pusblished to the DLQ `%s` (Id:`%s`): ```%s```",
			conf.SQS.DLQName, *message.MessageId, *message.Body,
		)
		return slack.PostWebhook(conf.Slack.WebhookUrl, &slack.WebhookMessage{Text: text})
	})

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
