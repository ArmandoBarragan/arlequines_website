package settings

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
)

type Config struct {
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBPassword                  string
	DBName                      string
	DBSSLMode                   string
	HostURL                     string
	StripePublicKey             string
	StripePrivateKey            string
	RedisPassword               string
	RedisConsumerConfigurations map[string]string
	SMTPHost                    string
	SMTPPort                    string
	SMTPUser                    string
	SMTPPass                    string
	ConsumerTimer               time.Duration
}

func InitRedis(config *Config) *redis.Client {
	stripe.Key = config.StripePrivateKey
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: config.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return client
}

func InitConsumerGroup(redisClient *redis.Client, config *Config) {
	ctx := context.Background()
	err := redisClient.XGroupCreateMkStream(
		ctx,
		config.RedisConsumerConfigurations["payment_stream_name"],
		config.RedisConsumerConfigurations["payment_consumer_group_name"],
		config.RedisConsumerConfigurations["payment_consumer_group_id"],
	).Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Could not create consumer group: %v", err)
	}
}

func LoadConfig() *Config {
	timer, err := strconv.Atoi(os.Getenv("CONSUMER_TIMER"))
	if err != nil {
		log.Printf("Failed to load consumer timer")
		timer = 60
	}
	config := &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSLMode:        os.Getenv("DB_SSLMODE"),
		HostURL:          os.Getenv("HOST_URL"),
		StripePublicKey:  os.Getenv("STRIPE_PUBLIC_KEY"),
		StripePrivateKey: os.Getenv("STRIPE_PRIVATE_KEY"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisConsumerConfigurations: map[string]string{
			"payment_stream_name":         "payment_stream",
			"payment_consumer_prefix":     "payment_consumer_",
			"payment_consumer_group_name": "payment_consumer_group",
			"payment_consumer_group_id":   "0",
		},
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUser:      os.Getenv("SMTP_USER"),
		SMTPPass:      os.Getenv("SMTP_PASS"),
		ConsumerTimer: time.Duration(timer),
	}
	return config
}
