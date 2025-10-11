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
		fmt.Sprintf("%s-mp-log-counter", p.Prefix),
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
	if !p.PerSec {
		return float64(i) / p.duration
	}
	return float64(i) / p.duration / 60
}

func (p *Opt) Output() string {
	now := uint64(time.Now().Unix())
	var output strings.Builder
	for _, af := range p.aggregatorFunctions {
		switch af.aggregator {
		case "count":
			fmt.Fprintf(&output, "%s.%s\t%f\t%d\n", p.Prefix, af.name, p.calculatePerDuration(af.count), now)
		case "group_by":
			for k, v := range af.groupBy {
				safeKey := strings.ReplaceAll(k, " ", "_")
				safeKey = strings.ReplaceAll(safeKey, ".", "_")
				fmt.Fprintf(&output, "%s.%s.%s\t%f\t%d\n", p.Prefix, af.name, safeKey, p.calculatePerDuration(v), now)
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
