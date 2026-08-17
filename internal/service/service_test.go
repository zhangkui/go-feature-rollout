package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func intPointer(value int) *int { return &value }

func TestFlagLifecycle(t *testing.T) {
	flagService := New()
	created, err := flagService.CreateFlag("checkout-redesign", "new checkout")
	if err != nil {
		t.Fatalf("CreateFlag returned error: %v", err)
	}
	if created.Key != "checkout-redesign" {
		t.Fatalf("unexpected key %q", created.Key)
	}
	if _, err := flagService.CreateFlag("checkout-redesign", "duplicate"); !errors.Is(err, ErrFlagExists) {
		t.Fatalf("expected ErrFlagExists, got %v", err)
	}
	if err := flagService.ArchiveFlag(created.Key); err != nil {
		t.Fatalf("ArchiveFlag returned error: %v", err)
	}
	archived, err := flagService.GetFlag(created.Key)
	if err != nil {
		t.Fatalf("GetFlag returned error: %v", err)
	}
	if !archived.Archived {
		t.Fatal("expected archived flag")
	}
}

func TestValidateRuleAcceptsSupportedConditions(t *testing.T) {
	start := time.Date(2026, time.August, 17, 9, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	rule := model.Rule{
		ID:       "employees",
		Priority: 10,
		Conditions: []model.Condition{
			{Attribute: "environment", Operator: model.OperatorEquals, Value: "production"},
			{Attribute: "team", Operator: model.OperatorIn, Values: []string{"payments", "search"}},
		},
		Percentage: intPointer(25),
		Enabled:    true,
		StartAt:    &start,
		EndAt:      &end,
	}
	if err := ValidateRule(rule); err != nil {
		t.Fatalf("ValidateRule returned error: %v", err)
	}
}

func TestPublishEvaluateAndTrace(t *testing.T) {
	flagService := New()
	if _, err := flagService.CreateFlag("search-v2", "new search"); err != nil {
		t.Fatal(err)
	}
	rule := model.Rule{
		ID:       "production-staff",
		Priority: 1,
		Conditions: []model.Condition{
			{Attribute: "environment", Operator: model.OperatorEquals, Value: "production"},
			{Attribute: "role", Operator: model.OperatorIn, Values: []string{"staff", "admin"}},
		},
		Percentage: intPointer(100),
		Enabled:    true,
	}
	if err := flagService.UpsertDraftRule("search-v2", rule); err != nil {
		t.Fatal(err)
	}
	published, err := flagService.Publish("search-v2", time.Date(2026, time.August, 17, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	decision, err := flagService.Evaluate("search-v2", published.Number, model.EvaluationRequest{
		SubjectID:  "user-42",
		Attributes: map[string]string{"environment": "production", "role": "staff"},
		At:         time.Date(2026, time.August, 17, 10, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !decision.Enabled || decision.MatchedRuleID != rule.ID {
		t.Fatalf("unexpected decision: %+v", decision)
	}
	if len(decision.Trace) != 1 || !decision.Trace[0].Matched {
		t.Fatalf("unexpected trace: %+v", decision.Trace)
	}
}

func TestRulePriorityUsesLowestNumberFirst(t *testing.T) {
	flagService := New()
	if _, err := flagService.CreateFlag("priority", "priority ordering"); err != nil {
		t.Fatal(err)
	}
	for _, rule := range []model.Rule{
		{ID: "later", Priority: 20, Enabled: true, Percentage: intPointer(100)},
		{ID: "first", Priority: 10, Enabled: true, Percentage: intPointer(100)},
	} {
		if err := flagService.UpsertDraftRule("priority", rule); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := flagService.Publish("priority", time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	decision, err := flagService.Evaluate("priority", 0, model.EvaluationRequest{SubjectID: "subject"})
	if err != nil {
		t.Fatal(err)
	}
	if decision.MatchedRuleID != "first" {
		t.Fatalf("expected first rule, got %q", decision.MatchedRuleID)
	}
}

func TestBatchEvaluateReturnsInputOrder(t *testing.T) {
	flagService := New()
	if _, err := flagService.CreateFlag("batch", "batch evaluation"); err != nil {
		t.Fatal(err)
	}
	if err := flagService.UpsertDraftRule("batch", model.Rule{ID: "all", Enabled: true, Percentage: intPointer(100)}); err != nil {
		t.Fatal(err)
	}
	if _, err := flagService.Publish("batch", time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	requests := []model.EvaluationRequest{{SubjectID: "one"}, {SubjectID: "two"}}
	results, err := flagService.BatchEvaluate(context.Background(), "batch", 0, requests)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != len(requests) {
		t.Fatalf("expected %d results, got %d", len(requests), len(results))
	}
	for index, result := range results {
		if !result.Enabled {
			t.Fatalf("result %d should be enabled", index)
		}
	}
}
