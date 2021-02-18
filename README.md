# DLQ X9

DLQ-X9 is a straightforward application that sends a message in a Slack channel
every time it detects a new message in a DLQ.

It's quite useful for debugging your application or messaging system.

**Note**: Every time a message is sent to Slack, it'll be deleted from the queue.

## Docker

You can run this project using the available image.

```bash
docker run -it \
  -e AWS_ACCESS_KEY_ID='<id>' \
  -e AWS_SECRET_ACCESS_KEY='<key>' \
  -e AWS_REGION='<region>' \
  dlq-x9 \
    --slack.webhook-url '<url>' \
    --sqs.dlq-name '<name>' \
    --internal.worker-pool '<size>'
```
