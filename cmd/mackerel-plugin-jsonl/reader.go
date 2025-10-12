package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kazeburo/followparser"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/montanaflynn/stats"
)

func (p *Opt) run() error {
	parser := NewParser(p)
	fp := &followparser.Parser{
		WorkDir:  pluginutil.PluginWorkDir(),
		Callback: parser,
		Silent:   !p.Verbose,
	}
	if p.LogArchiveDir != "" {
		fp.ArchiveDir = p.LogArchiveDir
	}
	_, err := fp.Parse(
		fmt.Sprintf("%s-mp-jsonl", p.Prefix),
		p.LogFile,
	)
	if err != nil {
		return err
	}
	output := p.Output()
	fmt.Print(output)
	return nil
}

func (p *Opt) calculatePerDuration(i int) float64 {
	if p.PerSec {
		return float64(i) / p.duration
	}
	return (float64(i) / p.duration) * 60
}

func (p *Opt) Output() string {
	now := uint64(time.Now().Unix())
	var output strings.Builder
	for i := 0; i < len(p.aggregatorFunctions); i++ {
		af := p.aggregatorFunctions[i]
		switch af.aggregator {
		case "count":
			if p.duration == 0 {
				// avoid division by zero
				continue
			}
			fmt.Fprintf(&output, "%s.%s\t%f\t%d\n", p.Prefix, af.name, p.calculatePerDuration(af.count), now)
		case "group_by", "group_by_with_percentage":
			if p.duration == 0 {
				// avoid division by zero
				continue
			}
			modifiedMap := map[string]int{}
			for k, v := range af.groupBy {
				safeKey := strings.ReplaceAll(k, " ", "_")
				safeKey = strings.ReplaceAll(safeKey, ".", "_")
				modifiedKey := af.applyModifiers(safeKey)
				modifiedMap[modifiedKey] += v
			}
			af.groupBy = modifiedMap
			total := 0
			for k, v := range af.groupBy {
				safeKey := strings.ReplaceAll(k, " ", "_")
				safeKey = strings.ReplaceAll(safeKey, ".", "_")
				fmt.Fprintf(&output, "%s.%s.%s\t%f\t%d\n", p.Prefix, af.name, safeKey, p.calculatePerDuration(v), now)
				total += v
			}
			if af.aggregator == "group_by_with_percentage" && total > 0 {
				for k, v := range af.groupBy {
					safeKey := strings.ReplaceAll(k, " ", "_")
					safeKey = strings.ReplaceAll(safeKey, ".", "_")
					percentage := float64(v) / float64(total) * 100
					fmt.Fprintf(&output, "%s.%s_percentage.%s\t%f\t%d\n", p.Prefix, af.name, safeKey, percentage, now)
				}
			}
		case "percentile":
			if len(af.percentiles) == 0 {
				continue
			}
			mean, _ := stats.Mean(af.percentiles)
			fmt.Fprintf(&output, "%s.%s.mean\t%f\t%d\n", p.Prefix, af.name, mean, now)
			ptiles := map[string]float64{
				"90": 90.0,
				"95": 95.0,
				"99": 99.0,
			}
			for name, ptile := range ptiles {
				ptile, err := stats.Percentile(af.percentiles, ptile)
				if err != nil {
					continue
				}
				fmt.Fprintf(&output, "%s.%s.p%s\t%f\t%d\n", p.Prefix, af.name, name, ptile, now)
			}
		}
	}
	return output.String()
}
