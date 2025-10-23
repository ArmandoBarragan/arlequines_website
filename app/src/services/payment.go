package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Config *settings.Config

type PaymentEvent struct {
	Email          string `redis:"email" json:"email"`
	Amount         int64  `redis:"amount" json:"amount"`
	Quantity       int64  `redis:"quantity" json:"quantity"`
	PresentationID uint   `redis:"presentation_id" json:"presentation_id"`
}

func (p *PaymentEvent) CreateEmailSendingEventRedis() {
	RedisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: Config.RedisConsumerConfigurations["payment_stream_name"],
		Values: map[string]any{
			"email":           p.Email,
			"amount":          p.Amount,
			"quantity":        p.Quantity,
			"presentation_id": p.PresentationID,
		},
	})
}

// New method to send to SQS
func (p *PaymentEvent) CreateEmailSendingEventSQS() error {
	// Create AWS config
	localConfig := settings.LoadConfig()
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(localConfig.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			localConfig.AWSAccessKeyID,
			localConfig.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(awsConfig)
	// Convert PaymentEvent to JSON
	messageBody, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payment event: %w", err)
	}

	// Send message to SQS
	_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(localConfig.SQSQueueURL),
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"EventType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("PaymentSuccess"),
			},
			"PresentationID": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprintf("%d", p.PresentationID)),
			},
		},
		MessageGroupId: aws.String("payment_group"),
		MessageDeduplicationId: aws.String(fmt.Sprintf("%s-%s%d", p.PresentationID, p.Email, rand.Intn(100))),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	fmt.Printf("Successfully sent payment event to SQS for presentation %d", p.PresentationID)
	return nil
}
