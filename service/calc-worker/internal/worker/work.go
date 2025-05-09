package worker

import (
	"context"
	"log"

	kaf "calc-worker/internal/kafka"

	"github.com/golkity/Calc_2.0/pkg/calc"
)

type Service struct{ prod *kaf.Producer }

func New(p *kaf.Producer) *Service { return &Service{p} }

func (s *Service) Handle(ctx context.Context, t kaf.Task) {
	res, err := calc.Calc(t.Expression)
	if err != nil {
		log.Printf("calc err: %v", err)
		return
	}
	if err := s.prod.Publish(ctx, t.ExprID, res); err != nil {
		log.Printf("publish err: %v", err)
	}
}
