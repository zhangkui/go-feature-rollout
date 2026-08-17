package service

import (
	"context"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func (s *Service) BatchEvaluate(ctx context.Context, key string, versionNumber int, requests []model.EvaluationRequest) ([]model.Decision, error) {
	// Honor a pre-canceled context before any evaluation work begins so no
	// partial decisions are produced. Mid-processing cancellation is still
	// handled inside the loop below, where the already-completed prefix is
	// returned alongside the cancellation error.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
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
