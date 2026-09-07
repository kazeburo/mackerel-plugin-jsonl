package main

import (
	"bytes"
	"log"

	"github.com/buger/jsonparser"
)

func (opt *Opt) jsonParsed(idx int, value []byte, vt jsonparser.ValueType, err error) {
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	if (vt == jsonparser.NotExist) || (vt == jsonparser.Null) {
		return
	}

	err = opt.aggregatorFunctions[idx].appendData(value)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
}

func (opt *Opt) Parse(b []byte) error {
	if opt.filterByte != nil && !bytes.Contains(b, *opt.filterByte) {
		return nil
	}
	if opt.ignoreByte != nil && bytes.Contains(b, *opt.ignoreByte) {
		return nil
	}
	if opt.SkipUntilBracket {
		i := bytes.IndexByte(b, '{')
		if i > 0 {
			b = b[i:]
		}
	}

	jsonparser.EachKey(b, opt.jsonParsed, opt.paths...)
	return nil
}

func (opt *Opt) Finish(duration float64) {
	opt.duration = duration
}
