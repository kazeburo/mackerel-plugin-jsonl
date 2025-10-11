package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
)

var version string

type AggregatorFunction struct {
	name        string
	jsonKey     []string
	aggregator  string
	count       int
	groupBy     map[string]int
	percentiles []float64
}

type Opt struct {
	Version             bool     `short:"v" long:"version" description:"Show version"`
	Filter              string   `long:"filter" description:"filter string used before check pattern."`
	Ignore              string   `long:"ignore" description:"ignore string used before check pattern."`
	KeyNames            []string `short:"k" long:"key-name" required:"true" description:"Key name for json path"`
	JsonKeys            []string `short:"j" long:"json-path" required:"true" description:"JSON key (dot concatenated) to extract log message."`
	Aggregator          []string `short:"a" long:"aggregator" required:"true" description:"Aggregator type. valid values are count, group_by, percentile. count is default." choice:"count" choice:"group_by" choice:"percentile"`
	SkipUntilBracket    bool     `long:"skip-until-json" description:"skip reading until first { for json log with plain text header"`
	Prefix              string   `long:"prefix" required:"true" description:"Metric key prefix"`
	PerSec              bool     `long:"per-second" description:"calculate per-seconds count. default per minute count"`
	LogFile             string   `short:"l" long:"log-file" description:"Path to log file" required:"true"`
	LogArchiveDir       string   `long:"log-archive-dir" default:"" description:"Path to log archive directory"`
	Verbose             bool     `long:"verbose" description:"display infomational logs"`
	aggregatorFunctions []AggregatorFunction
	filterByte          *[]byte
	ignoreByte          *[]byte
	duration            float64
}

func splitUnescapedDot(s string) []string {
	s = strings.TrimSpace(s)
	var result []string
	current := ""
	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && !escaped {
			escaped = true
			continue
		}
		if c == '.' && !escaped {
			result = append(result, current)
			current = ""
		} else {
			if escaped {
				current += "\\"
			}
			current += string(c)
		}
		escaped = false
	}
	result = append(result, current)
	return result
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
		keys := splitUnescapedDot(opt.JsonKeys[i])
		af := AggregatorFunction{
			name:       opt.KeyNames[i],
			jsonKey:    keys,
			aggregator: opt.Aggregator[i],
		}
		opt.aggregatorFunctions = append(opt.aggregatorFunctions, af)
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
