package service

import (
	"errors"
	"strings"
	"sync"

	"github.com/zhangkui/go-feature-rollout/internal/model"
)

var (
	ErrFlagNotFound    = errors.New("flag not found")
	ErrFlagExists      = errors.New("flag key already exists")
	ErrFlagArchived    = errors.New("flag is archived")
	ErrVersionNotFound = errors.New("version not found")
	ErrRuleExists      = errors.New("rule id already exists")
	ErrRuleNotFound    = errors.New("rule not found")
	ErrNoPublished     = errors.New("flag has no published version")
)

type Service struct {
	mu    sync.RWMutex
	flags map[string]*model.Flag
}

func New() *Service {
	return &Service{flags: make(map[string]*model.Flag)}
}

func canonicalKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}

func (s *Service) CreateFlag(key, description string) (model.Flag, error) {
	rawKey := strings.TrimSpace(key)
	if rawKey == "" {
		return model.Flag{}, errors.New("flag key is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.flags[rawKey]; exists {
		return model.Flag{}, ErrFlagExists
	}

	canonical := canonicalKey(rawKey)
	flag := &model.Flag{
		Key:         canonical,
		Description: strings.TrimSpace(description),
		Draft: model.Version{
			Default: false,
			Source:  "new",
		},
	}
	s.flags[canonical] = flag
	return cloneFlag(*flag), nil
}

func (s *Service) GetFlag(key string) (model.Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return model.Flag{}, ErrFlagNotFound
	}
	return cloneFlag(*flag), nil
}

func (s *Service) ArchiveFlag(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[canonicalKey(key)]
	if !ok {
		return ErrFlagNotFound
	}
	flag.Archived = true
	return nil
}

func cloneFlag(flag model.Flag) model.Flag {
	cloned := flag
	cloned.Draft = deepCloneVersion(flag.Draft)
	cloned.Versions = make([]model.Version, len(flag.Versions))
	for index := range flag.Versions {
		cloned.Versions[index] = deepCloneVersion(flag.Versions[index])
	}
	return cloned
}
