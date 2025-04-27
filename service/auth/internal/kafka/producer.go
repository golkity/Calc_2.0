package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const topicAuthEvents = "auth-events"

type Producer struct{ w *kafka.Writer }

func New(brokers []string) *Producer {
	return &Producer{
		w: &kafka.Writer{
			Addr:      kafka.TCP(brokers...),
			Topic:     topicAuthEvents,
			Balancer:  &kafka.Hash{},
			BatchSize: 1,
		},
	}
}

func (p *Producer) Close(ctx context.Context) error {
	return p.w.Close()
}

func (p *Producer) PublishUserCreated(ctx context.Context, id int64, email string) error {
	return p.publish(ctx, id, map[string]interface{}{
		"event": "user_created", "id": id, "email": email, "ts": time.Now(),
	})
}

func (p *Producer) PublishLogin(ctx context.Context, id int64) error {
	return p.publish(ctx, id, map[string]interface{}{
		"event": "login", "id": id, "ts": time.Now(),
	})
}

func (p *Producer) publish(ctx context.Context, id int64, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(strconv.FormatInt(id, 10)),
		Value: b,
	}

	if err = p.w.WriteMessages(ctx, msg); err != nil {
		log.Printf("kafka publish error: %v", err)
	}
	return err
}
