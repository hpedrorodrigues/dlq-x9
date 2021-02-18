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
}

func New(queueName string, workerPool int) *consumer {
	return &consumer{queueName, workerPool}
}

func (c *consumer) Consume(
	log *logrus.Logger,
	fn func(message *sqs.Message, log *logrus.Entry) error,
) {
	for w := 1; w <= c.workerPool; w++ {
		go c.worker(fn, log.WithField("worker", w))
	}
}

func (c *consumer) worker(
	fn func(message *sqs.Message, log *logrus.Entry) error,
	log *logrus.Entry,
) {
	log.Info("Starting")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(c.queueName),
	})
	if err != nil {
		log.Fatalf("Fatal error retrieving queue URL: %v", err)
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
			log.Errorf("Error retrieving messages: %v", err)
			continue
		}

		var wg sync.WaitGroup
		for _, message := range output.Messages {
			wg.Add(1)
			go func(m *sqs.Message) {
				defer wg.Done()
				if err := fn(m, log); err != nil {
					log.Errorf("Error processing message (%s): %v", *m.MessageId, err)
					return
				}

				_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      queueUrl,
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					log.Errorf("Error deleting message (%s): %v", *m.MessageId, err)
					return
				}
			}(message)

			wg.Wait()
		}
	}
}
