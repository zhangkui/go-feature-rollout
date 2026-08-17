package service

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func bucketFor(flagKey, subjectID string) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(flagKey + ":" + subjectID))
	return int(hasher.Sum32() % 100)
}

func (s *Service) Evaluate(key string, versionNumber int, request model.EvaluationRequest) (model.Decision, error) {
	if request.At.IsZero() {
		request.At = time.Now().UTC()
	}
	version, err := s.PublishedVersion(key, versionNumber)
	if err != nil {
		return model.Decision{}, err
	}
	decision := model.Decision{Enabled: version.Default, Version: version.Number}
	for _, rule := range sortedRules(version.Rules) {
		matched, reason := matchesRule(rule, request)
		step := model.TraceStep{RuleID: rule.ID, Matched: matched, Reason: reason}
		decision.Trace = append(decision.Trace, step)
		if !matched || !rule.Enabled {
			continue
		}
		if rule.Percentage != nil {
			bucket := bucketFor(key, request.SubjectID)
			decision.Bucket = &bucket
			if bucket >= *rule.Percentage {
				continue
			}
		}
		decision.Enabled = true
		decision.MatchedRuleID = rule.ID
		return decision, nil
	}
	return decision, nil
}

func matchesRule(rule model.Rule, request model.EvaluationRequest) (bool, string) {
	if rule.StartAt != nil && request.At.Before(*rule.StartAt) {
		return false, "before start time"
	}
	if rule.EndAt != nil && !request.At.Before(*rule.EndAt) {
		return false, "at or after end time"
	}
	for _, condition := range rule.Conditions {
		value := request.Attributes[condition.Attribute]
		switch condition.Operator {
		case model.OperatorEquals:
			if value != condition.Value {
				return false, "attribute does not equal expected value"
			}
		case model.OperatorNotEqual:
			if value == condition.Value {
				return false, "attribute equals excluded value"
			}
		case model.OperatorIn:
			if !contains(condition.Values, value) {
				return false, "attribute is not in allowed set"
			}
		case model.OperatorContains:
			if !strings.Contains(value, condition.Value) {
				return false, "attribute does not contain expected text"
			}
		default:
			return false, fmt.Sprintf("unsupported operator %q", condition.Operator)
		}
	}
	return true, "all conditions matched"
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func stableRuleIDs(rules []model.Rule) []string {
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	sort.Strings(ids)
	return ids
}
