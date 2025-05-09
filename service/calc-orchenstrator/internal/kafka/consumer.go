package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

const topicTasks = "calc.tasks"

type TaskReq struct {
	UserID     int64  `json:"user_id"`
	Expression string `json:"expression"`
}

func Consume(ctx context.Context, brokers []string, handle func(TaskReq)) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "calc-orchestrator",
		Topic:   topicTasks,
	})
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}
		var req TaskReq
		if err := json.Unmarshal(m.Value, &req); err != nil {
			log.Printf("orchestrator: bad msg %v", err)
			continue
		}
		handle(req)
	}
}
