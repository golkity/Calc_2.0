package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

const topicResults = "calc.results"

type Producer struct{ w *kafka.Writer }

func New(brokers []string) *Producer {
	return &Producer{w: &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topicResults,
		Balancer: &kafka.Hash{},
	}}
}

func (p *Producer) PublishResult(ctx context.Context, exprID int64, res float64) error {
	b, _ := json.Marshal(map[string]any{"expr_id": exprID, "result": res})
	return p.w.WriteMessages(ctx, kafka.Message{Value: b})
}

func (p *Producer) Close() error { return p.w.Close() }
