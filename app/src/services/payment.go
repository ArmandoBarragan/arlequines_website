package services

import (
	"context"

	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Config *settings.Config

type PaymentEvent struct {
	Email          string `redis:"email"`
	Amount         int64  `redis:"amount"`
	Quantity       int64  `redis:"quantity"`
	PresentationID uint   `redis:"presentation_id"`
}

func (p *PaymentEvent) CreateEmailSendingEvent() {
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
