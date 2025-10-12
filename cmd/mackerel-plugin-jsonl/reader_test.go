package main

import (
	"testing"
)

func TestOpt_calculatePerDuration(t *testing.T) {
	opt := &Opt{duration: 60, PerSec: false}
	if v := opt.calculatePerDuration(120); v != 2 {
		t.Errorf("expected 2, got %v", v)
	}
	opt.PerSec = true
	if v := opt.calculatePerDuration(120); v != 2/60.0 {
		t.Errorf("expected %v, got %v", 2/60.0, v)
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
	out := opt.Output()
	if out == "" {
		t.Errorf("expected output, got empty string")
	}
}
