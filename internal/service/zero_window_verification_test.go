package service

import (
	"testing"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func TestValidateRuleRejectsZeroLengthWindow(t *testing.T) {
	instant := time.Date(2026, time.August, 17, 12, 0, 0, 0, time.UTC)
	rule := model.Rule{ID: "instant", Priority: 1, Enabled: true, StartAt: &instant, EndAt: &instant}
	if err := ValidateRule(rule); err == nil {
		t.Fatal("expected equal start and end times to be rejected")
	}
}
