package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

const topicTasks = "calc.tasks"

type Task struct {
	ExprID     int64  `json:"expr_id"`
	Expression string `json:"expression"`
}

func Consume(ctx context.Context, brokers []string, handle func(Task)) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "calc-workers",
		Topic:   topicTasks,
	})
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var t Task
		if err := json.Unmarshal(m.Value, &t); err != nil {
			log.Printf("worker: bad msg %v", err)
			continue
		}
		handle(t)
	}
}
