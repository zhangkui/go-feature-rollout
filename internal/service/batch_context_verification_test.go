package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func TestBatchEvaluateHonorsPreCanceledContext(t *testing.T) {
	flagService := New()
	const key = "batch-cancel"
	if _, err := flagService.CreateFlag(key, "batch"); err != nil {
		t.Fatal(err)
	}
	if err := flagService.SetDraftDefault(key, true); err != nil {
		t.Fatal(err)
	}
	if _, err := flagService.Publish(key, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	results, err := flagService.BatchEvaluate(ctx, key, 0, []model.EvaluationRequest{{SubjectID: "first"}, {SubjectID: "second"}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("pre-canceled batch must not evaluate any items, got %d results", len(results))
	}
}
