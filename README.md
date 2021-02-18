# DLQ X9

DLQ-X9 is a straightforward application that sends a message in a Slack channel
every time it detects a new message in an SQS DLQ.

It's quite useful for debugging your application or messaging system.

**Note**: Every time a message is sent to Slack, it'll be deleted from the queue.

## Running

You can run this project using the available Docker image.

```bash
docker run -it \
  -e AWS_ACCESS_KEY_ID='<id>' \
  -e AWS_SECRET_ACCESS_KEY='<key>' \
  -e AWS_REGION='<region>' \
  ghcr.io/hpedrorodrigues/dlq-x9:0.1.0 \
    --slack.webhook-url '<url>' \
    --sqs.dlq-name '<name>' \
    --internal.worker-pool '<size>'
```

You can also run this project in your Kubernetes cluster using the manifest
files inside the `manifests` folder.

Edit the configuration files there, apply them and voilà.
