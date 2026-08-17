package model

import "time"

type Operator string

const (
	OperatorEquals   Operator = "eq"
	OperatorNotEqual Operator = "neq"
	OperatorIn       Operator = "in"
	OperatorContains Operator = "contains"
)

type Condition struct {
	Attribute string   `json:"attribute"`
	Operator  Operator `json:"operator"`
	Value     string   `json:"value,omitempty"`
	Values    []string `json:"values,omitempty"`
}

type Rule struct {
	ID         string      `json:"id"`
	Priority   int         `json:"priority"`
	Conditions []Condition `json:"conditions"`
	Percentage *int        `json:"percentage,omitempty"`
	Enabled    bool        `json:"enabled"`
	StartAt    *time.Time  `json:"start_at,omitempty"`
	EndAt      *time.Time  `json:"end_at,omitempty"`
}

type Version struct {
	Number      int        `json:"number"`
	Rules       []Rule     `json:"rules"`
	Default     bool       `json:"default"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Source      string     `json:"source,omitempty"`
}

type Flag struct {
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Archived    bool      `json:"archived"`
	Draft       Version   `json:"draft"`
	Versions    []Version `json:"versions"`
}

type EvaluationRequest struct {
	SubjectID  string            `json:"subject_id"`
	Attributes map[string]string `json:"attributes,omitempty"`
	At         time.Time         `json:"at,omitempty"`
}

type TraceStep struct {
	RuleID  string `json:"rule_id"`
	Matched bool   `json:"matched"`
	Reason  string `json:"reason"`
}

type Decision struct {
	Enabled       bool        `json:"enabled"`
	Version       int         `json:"version"`
	MatchedRuleID string      `json:"matched_rule_id,omitempty"`
	Bucket        *int        `json:"bucket,omitempty"`
	Trace         []TraceStep `json:"trace"`
}
