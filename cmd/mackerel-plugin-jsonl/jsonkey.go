package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-andiamo/splitter"
)

var ErrEmptyJsonKey = errors.New("json key is empty")

// path.to.key | toupper | tolower | trimspace | replace(regex, repl)
func parseJsonKeyWithFunc(s string) ([]string, []JsonKeyModifier, error) {
	emptyJsonModifier := []JsonKeyModifier{}
	// , splitter.Parenthesis, splitter.SquareBracket
	sp, _ := splitter.NewSplitter('|', splitter.DoubleQuotesBackSlashEscaped, splitter.SingleQuotesDoubleEscaped, splitter.Parenthesis, splitter.SquareBrackets)
	// do not unescapeQuotes. do split more times
	keys, err := sp.Split(s, splitter.TrimSpaces, splitter.IgnoreEmpties)
	if err != nil {
		return []string{}, emptyJsonModifier, err
	}
	if len(keys) == 0 {
		return []string{}, emptyJsonModifier, ErrEmptyJsonKey
	}
	jsonKey, err := parseJsonKey(keys[0])
	if err != nil {
		return []string{}, emptyJsonModifier, err
	}
	modifiers := []JsonKeyModifier{}
	for _, fn := range keys[1:] {
		if fn == "tolower" {
			modifiers = append(modifiers, func(s string) string {
				return strings.ToLower(s)
			})
		} else if fn == "toupper" {
			modifiers = append(modifiers, func(s string) string {
				return strings.ToUpper(s)
			})
		} else if fn == "trimspace" {
			modifiers = append(modifiers, func(s string) string {
				return strings.TrimSpace(s)
			})
		} else if strings.HasPrefix(fn, "replace(") && strings.HasSuffix(fn, ")") {
			inner := fn[8 : len(fn)-1]
			// replace("pattern","repl")
			s, _ := splitter.NewSplitter(',', splitter.DoubleQuotesBackSlashEscaped, splitter.SingleQuotesDoubleEscaped)
			// must unescapeQuotes
			parts, err := s.Split(inner, splitter.TrimSpaces, splitter.UnescapeQuotes, splitter.IgnoreEmpties)
			if err != nil {
				return []string{}, emptyJsonModifier, err
			}
			if len(parts) != 2 {
				return []string{}, emptyJsonModifier, fmt.Errorf("invalid replace() format: %s", fn)
			}
			pattern := parts[0]
			reg, err := regexp.Compile(pattern) // validate regexp
			if err != nil {
				return []string{}, emptyJsonModifier, fmt.Errorf("invalid regexp: %w in %s", err, fn)
			}
			repl := parts[1]
			modifiers = append(modifiers, func(s string) string {
				return reg.ReplaceAllString(s, repl)
			})
		} else {
			return []string{}, emptyJsonModifier, fmt.Errorf("unknown modifier: %s", fn)
		}
	}
	return jsonKey, modifiers, nil
}

// path.to."foo.baz".[0].key
func parseJsonKey(s string) ([]string, error) {
	sp, _ := splitter.NewSplitter('.', splitter.DoubleQuotesBackSlashEscaped, splitter.SingleQuotesDoubleEscaped, splitter.SquareBrackets)
	keys, err := sp.Split(s, splitter.TrimSpaces, splitter.UnescapeQuotes, splitter.IgnoreEmpties)
	if err != nil {
		return []string{}, err
	}
	return keys, nil
}
