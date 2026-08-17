package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func ValidateRule(rule model.Rule) error {
	if strings.TrimSpace(rule.ID) == "" {
		return errors.New("rule id is required")
	}
	if rule.Priority < 0 {
		return errors.New("priority must not be negative")
	}
	if rule.Percentage != nil && (*rule.Percentage < 0 || *rule.Percentage > 100) {
		return errors.New("percentage must be between 0 and 100")
	}
	if rule.StartAt != nil && rule.EndAt != nil && rule.EndAt.Before(*rule.StartAt) {
		return errors.New("end time must be after start time")
	}
	for index, condition := range rule.Conditions {
		if strings.TrimSpace(condition.Attribute) == "" {
			return fmt.Errorf("condition %d attribute is required", index)
		}
		switch condition.Operator {
		case model.OperatorEquals, model.OperatorNotEqual, model.OperatorContains:
			if condition.Value == "" {
				return fmt.Errorf("condition %d value is required", index)
			}
		case model.OperatorIn:
			if len(condition.Values) == 0 {
				return fmt.Errorf("condition %d values are required", index)
			}
		default:
			return fmt.Errorf("condition %d has unsupported operator", index)
		}
	}
	return nil
}

func (s *Service) UpsertDraftRule(key string, rule model.Rule) error {
	if err := ValidateRule(rule); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return ErrFlagNotFound
	}
	if flag.Archived {
		return ErrFlagArchived
	}
	for index := range flag.Draft.Rules {
		if flag.Draft.Rules[index].ID != rule.ID {
			continue
		}
		flag.Draft.Rules[index].Priority = rule.Priority
		flag.Draft.Rules[index].Enabled = rule.Enabled
		flag.Draft.Rules[index].Percentage = cloneInt(rule.Percentage)
		flag.Draft.Rules[index].StartAt = cloneTime(rule.StartAt)
		flag.Draft.Rules[index].EndAt = cloneTime(rule.EndAt)
		flag.Draft.Rules[index].Conditions = append(flag.Draft.Rules[index].Conditions[:0], cloneConditions(rule.Conditions)...)
		return nil
	}
	flag.Draft.Rules = append(flag.Draft.Rules, deepCloneRule(rule))
	return nil
}

func (s *Service) SetDraftDefault(key string, enabled bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return ErrFlagNotFound
	}
	if flag.Archived {
		return ErrFlagArchived
	}
	flag.Draft.Default = enabled
	return nil
}

func cloneConditions(conditions []model.Condition) []model.Condition {
	cloned := make([]model.Condition, len(conditions))
	for index, condition := range conditions {
		cloned[index] = condition
		cloned[index].Values = append([]string(nil), condition.Values...)
	}
	return cloned
}
