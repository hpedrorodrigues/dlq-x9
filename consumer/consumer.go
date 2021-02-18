package consumer

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
  "github.com/sirupsen/logrus"
  "sync"
)

type consumer struct {
  queueName  string
  workerPool int
  log        *logrus.Logger
}

func New(queueName string, workerPool int, log *logrus.Logger) *consumer {
  return &consumer{queueName, workerPool, log}
}

func (c *consumer) Consume(fn func(message *sqs.Message) error) {
  for w := 1; w <= c.workerPool; w++ {
    go c.worker(w, fn)
  }
}

func (c *consumer) worker(id int, fn func(message *sqs.Message) error) {
  c.log.Infof("Starting worker: %d\n", id)

  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  svc := sqs.New(sess)

  result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
    QueueName: aws.String(c.queueName),
  })
  if err != nil {
    c.log.Fatalf("Fatal error retrieving queue URL: %v\n", err)
  }

  queueUrl := result.QueueUrl

  for {
    output, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
      AttributeNames: []*string{
        aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
      },
      MessageAttributeNames: []*string{
        aws.String(sqs.QueueAttributeNameAll),
      },
      QueueUrl:            queueUrl,
      MaxNumberOfMessages: aws.Int64(1),
    })
    if err != nil {
      c.log.Errorf("Error retrieving messages: %v\n", err)
      continue
    }

    var wg sync.WaitGroup
    for _, message := range output.Messages {
      wg.Add(1)
      go func(m *sqs.Message) {
        defer wg.Done()
        if err := fn(m); err != nil {
          return
        }

        _, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
          QueueUrl:      queueUrl,
          ReceiptHandle: message.ReceiptHandle,
        })
        if err != nil {
          c.log.Errorf("Error deleting message (%s): %v\n", *m.MessageId, err)
          return
        }
      }(message)

      wg.Wait()
    }
  }
}
