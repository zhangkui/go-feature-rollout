package service

import (
	"context"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func (s *Service) BatchEvaluate(ctx context.Context, key string, versionNumber int, requests []model.EvaluationRequest) ([]model.Decision, error) {
	results := make([]model.Decision, 0, len(requests))
	for _, request := range requests {
		decision, err := s.Evaluate(key, versionNumber, request)
		if err != nil {
			return nil, err
		}
		results = append(results, decision)
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}
	}
	return results, nil
}
