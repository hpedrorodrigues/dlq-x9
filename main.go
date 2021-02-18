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
  var log = logrus.New()

  conf := config.LoadConfiguration(log)
  cons := consumer.New(conf.SQS.DLQName, conf.Internal.WorkerPool, log)

  log.Info("DLQ-X9")

  cons.Consume(func(message *sqs.Message) error {
    text := fmt.Sprintf(
      "Hey, Amazon SQS moved a new message to the DLQ `%s` (Id:`%s`): ```%s```",
      conf.SQS.DLQName, *message.MessageId, *message.Body,
    )
    return slack.PostWebhook(conf.Slack.WebhookUrl, &slack.WebhookMessage{Text: text})
  })

  sigterm := make(chan os.Signal, 1)
  signal.Notify(sigterm, syscall.SIGTERM)
  signal.Notify(sigterm, syscall.SIGINT)
  <-sigterm
}
