package service

import "github.com/zhangkui/go-feature-rollout/internal/model"

func (s *Service) UpdateRuleConditions(key, ruleID string, conditions []model.Condition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return ErrFlagNotFound
	}
	for index := range flag.Draft.Rules {
		if flag.Draft.Rules[index].ID == ruleID {
			flag.Draft.Rules[index].Conditions = cloneConditions(conditions)
			return nil
		}
	}
	return ErrRuleNotFound
}
