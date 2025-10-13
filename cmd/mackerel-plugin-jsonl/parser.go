package main

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"unsafe"

	"github.com/buger/jsonparser"
)

type Parser struct {
	opt   *Opt
	paths [][]string
}

func NewParser(opt *Opt) *Parser {
	// initialize groupBy map and percentiles slice
	for i := range opt.aggregatorFunctions {
		if strings.HasPrefix(opt.aggregatorFunctions[i].aggregator, "group_by") {
			opt.aggregatorFunctions[i].groupBy = make(map[string]int)
		} else if opt.aggregatorFunctions[i].aggregator == "percentile" {
			opt.aggregatorFunctions[i].percentiles = []float64{}
		}
	}
	paths := [][]string{}
	for _, af := range opt.aggregatorFunctions {
		paths = append(paths, af.jsonKey)
	}
	return &Parser{
		opt:   opt,
		paths: paths,
	}
}

func bfloat64(b []byte) (float64, error) {
	return strconv.ParseFloat(unsafe.String(unsafe.SliceData(b), len(b)), 64)
}

func (p *Parser) jsonParsed(idx int, value []byte, vt jsonparser.ValueType, err error) {
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	if (vt == jsonparser.NotExist) || (vt == jsonparser.Null) {
		return
	}

	switch p.opt.aggregatorFunctions[idx].aggregator {
	case "count":
		p.opt.aggregatorFunctions[idx].count++
	case "group_by", "group_by_with_percentage":
		p.opt.aggregatorFunctions[idx].groupBy[string(value)]++
	case "percentile":
		floatValue, err := bfloat64(value)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		p.opt.aggregatorFunctions[idx].percentiles = append(p.opt.aggregatorFunctions[idx].percentiles, floatValue)
	}
}

func (p *Parser) Parse(b []byte) error {
	if p.opt.filterByte != nil && !bytes.Contains(b, *p.opt.filterByte) {
		return nil
	}
	if p.opt.ignoreByte != nil && bytes.Contains(b, *p.opt.ignoreByte) {
		return nil
	}
	if p.opt.SkipUntilBracket {
		i := bytes.IndexByte(b, '{')
		if i > 0 {
			b = b[i:]
		}
	}

	jsonparser.EachKey(b, p.jsonParsed, p.paths...)
	return nil
}

func (p *Parser) Finish(duration float64) {
	p.opt.duration = duration
}
