# sample-producer-consumer-app

A simple message processing system used as as a sample to improve using **AI**

## Architecture

```
┌──────────┐       ┌─────────┐       ┌──────────┐
│ Producer │──────▶│  SQS    │──────▶│ Consumer │
│  (Go)    │ send  │  Queue  │ read  │  (Go)    │
└──────────┘       └─────────┘       └────┬─────┘
                                          │ store
                                     ┌────▼─────┐
                                     │ DynamoDB  │
                                     │  Table    │
                                     └──────────┘
```

- **Producer** — Go app that sends a message to SQS every 2 seconds
- **Consumer** — Go app that reads messages from SQS and stores them in DynamoDB
- **Infra** — Terraform that deploys the SQS queue, DynamoDB table, and IAM role in AWS

## Project Structure

```
sample-producer-consumer-app/
├── producer/        # Go app
├── consumer/        # Go app
├── scripts/         # Start scripts
├── infra/           # Terraform (SQS, DynamoDB, IAM)
├── .env             # Environment variables
└── README.md
```

## Environment Variables

| Variable | Used by | Description |
|----------|---------|-------------|
| `AWS_REGION` | both | AWS region (default: `us-east-1`) |
| `SQS_QUEUE_URL` | both | SQS queue URL from Terraform output |
| `DYNAMODB_TABLE` | consumer | DynamoDB table name from Terraform output |

## Quick Start

**1. Deploy infrastructure:**

```bash
cd infra
terraform init
terraform apply
```

**2. Run the apps:**

```bash
# terminal 1
./scripts/start-producer.sh

# terminal 2
./scripts/start-consumer.sh
```
