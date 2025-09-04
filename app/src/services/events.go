package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/redis/go-redis/v9"
)
// Deprecated: Will use AWS SQS instead
func EmailEventConsumerWorker(redisClient *redis.Client, workerID int, config *settings.Config) {
	ctx := context.Background()
	myConsumerName := fmt.Sprintf("%s%d", config.RedisConsumerConfigurations["payment_consumer_prefix"], workerID)
	fmt.Printf("Starting Payment Consumer Worker '%s'...\n", myConsumerName)

	// Ensure the consumer group exists. Create it if it doesn't.
	// The "0" means start reading from the beginning of the stream if the group is new.
	err := redisClient.XGroupCreateConsumer(ctx, config.RedisConsumerConfigurations["payment_stream_name"], config.RedisConsumerConfigurations["payment_consumer_group_name"], myConsumerName).Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("FATAL: Worker '%s' failed to create consumer group: %v", myConsumerName, err)
	} else if err == nil {
		fmt.Printf("Worker '%s' created consumer group '%s' for stream '%s'.\n", myConsumerName, config.RedisConsumerConfigurations["payment_consumer_group_name"], config.RedisConsumerConfigurations["payment_stream_name"])
	} else {
		fmt.Printf("Worker '%s' joining existing consumer group '%s' for stream '%s'.\n", myConsumerName, config.RedisConsumerConfigurations["payment_consumer_group_name"], config.RedisConsumerConfigurations["payment_stream_name"])
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
		msg, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    config.RedisConsumerConfigurations["payment_consumer_group_name"],
			Consumer: myConsumerName,
			Streams:  []string{config.RedisConsumerConfigurations["payment_stream_name"], ">"},
			Count:    1,
			Block:    1 * time.Second,
		}).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Print(err) // TODO: Improve error handling
			continue
		}
		for _, stream := range msg {
			for _, message := range stream.Messages {
				fmt.Println(message.Values) // Add handler function
				err := redisClient.XAck(
					ctx,
					stream.Stream,
					config.RedisConsumerConfigurations["payment_consumer_group_name"],
					message.ID,
				).Err()
				if err != nil {
					log.Printf("Failed to acknowledge message %s: %v", message.ID, err)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
