package main

import (
	"strings"
	"testing"
)

func TestOpt_check(t *testing.T) {
	// 必須パラメータ不足
	opt := &Opt{}
	if err := opt.check(); err == nil {
		t.Error("expected error for missing params")
	}

	// パラメータ数不一致
	opt = &Opt{
		KeyNames:   []string{"foo"},
		JsonKeys:   []string{"foo"},
		Aggregator: []string{"count", "group_by"},
	}
	if err := opt.check(); err == nil {
		t.Error("expected error for param count mismatch")
	}

	// 不正なaggregator
	opt = &Opt{
		KeyNames:   []string{"foo"},
		JsonKeys:   []string{"foo"},
		Aggregator: []string{"invalid"},
	}
	if err := opt.check(); err == nil {
		t.Error("expected error for invalid aggregator")
	}

	// countでmodifier指定時はエラー
	opt = &Opt{
		KeyNames:   []string{"foo"},
		JsonKeys:   []string{"foo|tolower"},
		Aggregator: []string{"count"},
	}
	if err := opt.check(); err == nil {
		t.Error("expected error for modifier with count")
	}

	// percentileでmodifier指定時はエラー
	opt = &Opt{
		KeyNames:   []string{"foo"},
		JsonKeys:   []string{"foo|tolower"},
		Aggregator: []string{"percentile"},
	}
	if err := opt.check(); err == nil {
		t.Error("expected error for modifier with percentile")
	}

	// group_byでmodifier指定時はOK
	opt = &Opt{
		KeyNames:   []string{"foo"},
		JsonKeys:   []string{"foo|tolower"},
		Aggregator: []string{"group_by"},
	}
	if err := opt.check(); err != nil {
		t.Errorf("unexpected error for group_by with modifier: %v", err)
	}

	// 正常系
	opt = &Opt{
		KeyNames:   []string{"foo"},
		JsonKeys:   []string{"foo"},
		Aggregator: []string{"count"},
	}
	if err := opt.check(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestOpt_calculatePerDuration(t *testing.T) {
	opt := &Opt{duration: 60, PerSec: false}
	if v := opt.calculatePerDuration(120); v != 120 {
		t.Errorf("expected %v, got %v", 120, v)
	}
	opt.PerSec = true
	if v := opt.calculatePerDuration(120); v != 2 {
		t.Errorf("expected %v, got %v", 2, v)
	}
}

func TestOpt_Output(t *testing.T) {
	opt := &Opt{
		Prefix: "test",
		aggregatorFunctions: []*AggregatorFunction{
			{
				aggregator: "count",
				name:       "foo",
				count:      10,
			},
		},
		duration: 10,
	}
	out := opt.output()
	if out == "" {
		t.Errorf("expected output, got empty string")
	}
	if !strings.Contains(out, "test.foo\t60.000000\t") {
		t.Errorf("unexpected output: %s", out)
	}
}
