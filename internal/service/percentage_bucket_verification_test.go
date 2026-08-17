package service

import (
	"fmt"
	"hash/fnv"
	"testing"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func referenceBucket(flagKey, subjectID string) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(flagKey + ":" + subjectID))
	return int(hasher.Sum32() % 100)
}

func subjectForBucket(flagKey string, target int) string {
	for candidate := 0; candidate < 100000; candidate++ {
		subject := fmt.Sprintf("subject-%d", candidate)
		if referenceBucket(flagKey, subject) == target {
			return subject
		}
	}
	return ""
}

func TestEvaluatePercentageUsesHundredBuckets(t *testing.T) {
	flagService := New()
	const key = "gradual-rollout"
	if _, err := flagService.CreateFlag(key, "rollout"); err != nil {
		t.Fatal(err)
	}
	percentage := 99
	if err := flagService.UpsertDraftRule(key, model.Rule{ID: "ninety-nine", Priority: 1, Enabled: true, Percentage: &percentage}); err != nil {
		t.Fatal(err)
	}
	if _, err := flagService.Publish(key, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	subject := subjectForBucket(key, 99)
	if subject == "" {
		t.Fatal("failed to find deterministic bucket fixture")
	}
	decision, err := flagService.Evaluate(key, 0, model.EvaluationRequest{SubjectID: subject})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Enabled {
		t.Fatalf("subject in bucket 99 must be excluded from a 99%% rollout: %+v", decision)
	}
	if decision.Bucket == nil || *decision.Bucket != 99 {
		t.Fatalf("expected reported bucket 99, got %+v", decision.Bucket)
	}
}
