package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

const topicResults = "calc.results"

type Producer struct{ w *kafka.Writer }

func NewProducer(brokers []string) *Producer {
	return &Producer{w: &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topicResults,
		Balancer: &kafka.Hash{},
	}}
}

func (p *Producer) Publish(ctx context.Context, exprID int64, res float64) error {
	b, _ := json.Marshal(map[string]any{"expr_id": exprID, "result": res})
	return p.w.WriteMessages(ctx, kafka.Message{Value: b})
}

func (p *Producer) Ping(ctx context.Context) error {
	return p.w.WriteMessages(ctx, kafka.Message{Value: []byte(`{}`)})
}

func (p *Producer) Close() error { return p.w.Close() }
