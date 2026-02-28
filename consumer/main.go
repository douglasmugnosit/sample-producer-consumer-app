package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL environment variable is required")
	}

	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		log.Fatal("DYNAMODB_TABLE environment variable is required")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	ddbClient := dynamodb.NewFromConfig(cfg)

	log.Printf("consumer started | queue=%s table=%s region=%s", queueURL, tableName, region)

	for {
		result, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
		})
		if err != nil {
			log.Printf("ERROR receiving messages: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, msg := range result.Messages {
			messageID := aws.ToString(msg.MessageId)
			body := aws.ToString(msg.Body)

			_, err := ddbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item: map[string]types.AttributeValue{
					"id":           &types.AttributeValueMemberS{Value: messageID},
					"body":         &types.AttributeValueMemberS{Value: body},
					"processed_at": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
				},
			})
			if err != nil {
				log.Printf("ERROR storing message %s: %v", messageID, err)
				continue
			}

			_, err = sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("ERROR deleting message %s from queue: %v", messageID, err)
				continue
			}

			log.Printf("processed message %s: %s", messageID, body)
		}

		if len(result.Messages) == 0 {
			fmt.Print(".")
		}
	}
}
