package service

import (
	"testing"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func TestRollbackDraftDoesNotMutatePublishedVersion(t *testing.T) {
	flagService := New()
	const key = "immutable-history"
	if _, err := flagService.CreateFlag(key, "history"); err != nil {
		t.Fatal(err)
	}
	original := model.Rule{ID: "environment", Priority: 1, Enabled: true, Conditions: []model.Condition{{Attribute: "environment", Operator: model.OperatorEquals, Value: "production"}}}
	if err := flagService.UpsertDraftRule(key, original); err != nil {
		t.Fatal(err)
	}
	if _, err := flagService.Publish(key, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	if err := flagService.Rollback(key, 1); err != nil {
		t.Fatal(err)
	}
	changed := original
	changed.Conditions = []model.Condition{{Attribute: "environment", Operator: model.OperatorEquals, Value: "staging"}}
	if err := flagService.UpsertDraftRule(key, changed); err != nil {
		t.Fatal(err)
	}
	published, err := flagService.PublishedVersion(key, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := published.Rules[0].Conditions[0].Value; got != "production" {
		t.Fatalf("published version changed through rollback draft: got %q", got)
	}
}
