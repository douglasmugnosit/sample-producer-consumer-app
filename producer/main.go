package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL environment variable is required")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	counter := 0

	log.Printf("producer started | queue=%s region=%s", queueURL, region)

	for {
		counter++
		message := fmt.Sprintf("order-%d | timestamp=%s", counter, time.Now().UTC().Format(time.RFC3339))

		_, err := client.SendMessage(context.TODO(), &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(message),
		})
		if err != nil {
			log.Printf("ERROR sending message #%d: %v", counter, err)
		} else {
			log.Printf("sent message #%d: %s", counter, message)
		}

		time.Sleep(5 * time.Second)
	}
}
