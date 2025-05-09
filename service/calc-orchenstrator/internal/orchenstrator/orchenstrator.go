package orchestrator

import (
	"context"
	"log"

	kaf "github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/kafka"
	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/internal/repository"

	"github.com/golkity/Calc_2.0/pkg/calc"
)

type Service struct {
	exprRepo *repository.ExpressionRepo
	prod     *kaf.Producer
}

func New(er *repository.ExpressionRepo, prod *kaf.Producer) *Service {
	return &Service{er, prod}
}

func (s *Service) Handle(ctx context.Context, userID int64, expr string) {
	id, err := s.exprRepo.Create(ctx, userID, expr)
	if err != nil {
		log.Printf("expr create: %v", err)
		return
	}

	res, err := calc.Calc(expr)
	if err != nil {
		log.Printf("calc error: %v", err)
		return
	}
	if err := s.exprRepo.Done(ctx, id, res); err != nil {
		log.Printf("expr update: %v", err)
	}
	_ = s.prod.PublishResult(ctx, id, res)
}
