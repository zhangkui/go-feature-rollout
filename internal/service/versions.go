package service

import (
	"errors"
	"sort"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

func (s *Service) Publish(key string, now time.Time) (model.Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return model.Version{}, ErrFlagNotFound
	}
	if flag.Archived {
		return model.Version{}, ErrFlagArchived
	}
	for _, rule := range flag.Draft.Rules {
		if err := ValidateRule(rule); err != nil {
			return model.Version{}, err
		}
	}
	number := len(flag.Versions) + 1
	published := deepCloneVersion(flag.Draft)
	published.Number = number
	published.PublishedAt = cloneTime(&now)
	published.Source = "publish"
	flag.Versions = append(flag.Versions, published)
	flag.Draft = deepCloneVersion(published)
	flag.Draft.Number = number
	flag.Draft.PublishedAt = nil
	flag.Draft.Source = "draft"
	return deepCloneVersion(published), nil
}

func (s *Service) Rollback(key string, number int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return ErrFlagNotFound
	}
	if flag.Archived {
		return ErrFlagArchived
	}
	for _, version := range flag.Versions {
		if version.Number != number {
			continue
		}
		flag.Draft = deepCloneVersion(version)
		flag.Draft.Number = number
		flag.Draft.PublishedAt = nil
		flag.Draft.Source = "rollback"
		return nil
	}
	return ErrVersionNotFound
}

func (s *Service) PublishedVersion(key string, number int) (model.Version, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return model.Version{}, ErrFlagNotFound
	}
	if number == 0 {
		if len(flag.Versions) == 0 {
			return model.Version{}, ErrNoPublished
		}
		return deepCloneVersion(flag.Versions[len(flag.Versions)-1]), nil
	}
	for _, version := range flag.Versions {
		if version.Number == number {
			return deepCloneVersion(version), nil
		}
	}
	return model.Version{}, ErrVersionNotFound
}

func deepCloneVersion(version model.Version) model.Version {
	cloned := version
	cloned.Rules = make([]model.Rule, len(version.Rules))
	for index, rule := range version.Rules {
		cloned.Rules[index] = deepCloneRule(rule)
	}
	return cloned
}

func deepCloneRule(rule model.Rule) model.Rule {
	cloned := rule
	cloned.Conditions = cloneConditions(rule.Conditions)
	cloned.Percentage = cloneInt(rule.Percentage)
	cloned.StartAt = cloneTime(rule.StartAt)
	cloned.EndAt = cloneTime(rule.EndAt)
	return cloned
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func sortedRules(rules []model.Rule) []model.Rule {
	cloned := make([]model.Rule, len(rules))
	for index, rule := range rules {
		cloned[index] = deepCloneRule(rule)
	}
	sort.SliceStable(cloned, func(left, right int) bool {
		if cloned[left].Priority == cloned[right].Priority {
			return cloned[left].ID < cloned[right].ID
		}
		return cloned[left].Priority < cloned[right].Priority
	})
	return cloned
}

var errInvalidEvaluationTime = errors.New("evaluation time is required")
