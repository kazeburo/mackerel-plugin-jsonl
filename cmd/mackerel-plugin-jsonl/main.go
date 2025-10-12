package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
)

var version string

type JsonKeyModifier func(string) string

type AggregatorFunction struct {
	name             string
	jsonKey          []string
	JsonKeyModifiers []JsonKeyModifier
	aggregator       string
	count            int
	groupBy          map[string]int
	percentiles      []float64
}

func (af *AggregatorFunction) applyModifiers(s string) string {
	for _, mod := range af.JsonKeyModifiers {
		s = mod(s)
	}
	return s
}

type Opt struct {
	Version             bool     `short:"v" long:"version" description:"Show version"`
	Filter              string   `long:"filter" description:"filter string used before check pattern."`
	Ignore              string   `long:"ignore" description:"ignore string used before check pattern."`
	KeyNames            []string `short:"k" long:"key-name" required:"true" description:"Key name for json path"`
	JsonKeys            []string `short:"j" long:"json-key" required:"true" description:"JSON key and modifier functions to extract log message."`
	Aggregator          []string `short:"a" long:"aggregator" required:"true" description:"Aggregator type. valid values are count, group_by, group_by_with_percentage, percentile. count is default." choice:"count" choice:"group_by" choice:"group_by_with_percentage" choice:"percentile"`
	SkipUntilBracket    bool     `long:"skip-until-json" description:"skip reading until first { for json log with plain text header"`
	Prefix              string   `long:"prefix" required:"true" description:"Metric key prefix"`
	PerSec              bool     `long:"per-second" description:"calculate per-seconds count. default per minute count"`
	LogFile             string   `short:"l" long:"log-file" description:"Path to log file" required:"true"`
	LogArchiveDir       string   `long:"log-archive-dir" default:"" description:"Path to log archive directory"`
	Verbose             bool     `long:"verbose" description:"display infomational logs"`
	aggregatorFunctions []*AggregatorFunction
	filterByte          *[]byte
	ignoreByte          *[]byte
	duration            float64
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := Opt{}
	psr := flags.NewParser(&opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		fmt.Printf(`%s %s
Compiler: %s %s
`,
			os.Args[0],
			version,
			runtime.Compiler,
			runtime.Version())
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	if len(opt.KeyNames) == 0 {
		fmt.Fprint(os.Stderr, "Specify --key-name <name> --json-path <path> --aggregator <type>\n")
		return 1
	}

	if len(opt.KeyNames) != len(opt.JsonKeys) || len(opt.KeyNames) != len(opt.Aggregator) {
		fmt.Fprint(os.Stderr, "--key-name, --json-path and --aggregator must be specified the same number of times\n")
		return 1
	}

	for i := 0; i < len(opt.KeyNames); i++ {
		var keys []string
		var modifiers []JsonKeyModifier
		switch opt.Aggregator[i] {
		case "count", "percentile":
			keys, modifiers, err = parseJsonKeyWithFunc(opt.JsonKeys[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid json key: %v\n", err)
				return 1
			}
			if len(modifiers) > 0 {
				fmt.Fprintf(os.Stderr, "modifiers are not supported for %s aggregator\n", opt.Aggregator[i])
				return 1
			}
		case "group_by", "group_by_with_percentage":
			keys, modifiers, err = parseJsonKeyWithFunc(opt.JsonKeys[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid json key: %v\n", err)
				return 1
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown aggregator: %s\n", opt.Aggregator[i])
			return 1
		}
		af := AggregatorFunction{
			name:             opt.KeyNames[i],
			jsonKey:          keys,
			JsonKeyModifiers: modifiers,
			aggregator:       opt.Aggregator[i],
		}
		opt.aggregatorFunctions = append(opt.aggregatorFunctions, &af)
	}

	if opt.Filter != "" {
		b := []byte(opt.Filter)
		opt.filterByte = &b
	}
	if opt.Ignore != "" {
		b := []byte(opt.Ignore)
		opt.ignoreByte = &b
	}

	err = opt.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}
