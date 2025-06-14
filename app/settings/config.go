package settings

import (
	"os"

	"github.com/stripe/stripe-go/v82"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	HostURL          string
	StripePublicKey  string
	StripePrivateKey string
}

func LoadConfig() *Config {
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
	}
	stripe.Key = config.StripePrivateKey
	return config
}
