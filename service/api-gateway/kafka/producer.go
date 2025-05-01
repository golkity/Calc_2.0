package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/segmentio/kafka-go"
)

const topicTasks = "calc.tasks"

type Producer struct{ w *kafka.Writer }

func New(brokers []string) *Producer {
	return &Producer{w: &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topicTasks,
		Balancer: &kafka.Hash{},
	}}
}

func (p *Producer) PublishCalcRequest(ctx context.Context, userID int64, expr string) error {
	msg := map[string]any{"user_id": userID, "expression": expr}
	b, _ := json.Marshal(msg)
	return p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(userID, 10)),
		Value: b,
	})
}

func (p *Producer) Ping(ctx context.Context) error {
	return p.w.WriteMessages(ctx, kafka.Message{Value: []byte(`{}`)})
}

func (p *Producer) Close() error { return p.w.Close() }
