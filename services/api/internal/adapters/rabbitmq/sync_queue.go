package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	app "github.com/sutad-p/oh-my-trading/services/api/internal/application/marketdata"
)

type SyncPublisher struct {
	channel *amqp.Channel
	queue   string
}

func NewSyncPublisher(channel *amqp.Channel, queue string) *SyncPublisher {
	return &SyncPublisher{
		channel: channel,
		queue:   queue,
	}
}

func (p *SyncPublisher) PublishSyncRequest(ctx context.Context, request app.SyncRequest) error {
	if _, err := p.channel.QueueDeclare(p.queue, true, false, false, false, nil); err != nil {
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(ctx, "", p.queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}
