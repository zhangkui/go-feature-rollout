package service

import (
	"errors"
	"testing"
)

func TestCreateFlagRejectsCanonicalDuplicate(t *testing.T) {
	flagService := New()
	first, err := flagService.CreateFlag("search-v2", "first definition")
	if err != nil {
		t.Fatalf("first CreateFlag returned error: %v", err)
	}
	if _, err := flagService.CreateFlag(" SEARCH-V2 ", "replacement"); !errors.Is(err, ErrFlagExists) {
		t.Fatalf("expected ErrFlagExists for canonical duplicate, got %v", err)
	}
	stored, err := flagService.GetFlag("SEARCH-V2")
	if err != nil {
		t.Fatalf("GetFlag returned error: %v", err)
	}
	if stored.Description != first.Description {
		t.Fatalf("duplicate request replaced existing flag: got %q, want %q", stored.Description, first.Description)
	}
}
