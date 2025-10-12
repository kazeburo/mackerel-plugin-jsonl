package main

import (
	"reflect"
	"testing"
)

func TestParseJsonKeyWithFunc_basic(t *testing.T) {
	key, mods, err := parseJsonKeyWithFunc("foo.bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(key, []string{"foo", "bar"}) {
		t.Errorf("expected [foo bar], got %v", key)
	}
	if len(mods) != 0 {
		t.Errorf("expected no modifiers, got %d", len(mods))
	}
}

func TestParseJsonKeyWithFunc_modifiers(t *testing.T) {
	key, mods, err := parseJsonKeyWithFunc("foo|tolower|toupper|trimspace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key[0] != "foo" {
		t.Errorf("expected key 'foo', got %v", key)
	}
	if len(mods) != 3 {
		t.Errorf("expected 3 modifiers, got %d", len(mods))
	}
	// test modifier chain
	v := " FOO "
	for _, m := range mods {
		v = m(v)
	}
	if v != "FOO" {
		t.Errorf("modifier chain failed, got %v", v)
	}
}

func TestParseJsonKeyWithFunc_replace(t *testing.T) {
	{
		key, mods, err := parseJsonKeyWithFunc("foo|replace('o','a')")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key[0] != "foo" {
			t.Errorf("expected key 'foo', got %v", key)
		}
		if len(mods) != 1 {
			t.Errorf("expected 1 modifier, got %d", len(mods))
		}
		v := mods[0]("foo")
		if v != "faa" {
			t.Errorf("replace modifier failed, got %v", v)
		}
	}
	{
		key, mods, err := parseJsonKeyWithFunc(`foo|replace('[o|f]','a"')`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key[0] != "foo" {
			t.Errorf("expected key 'foo', got %v", key)
		}
		if len(mods) != 1 {
			t.Errorf("expected 1 modifier, got %d", len(mods))
		}
		v := mods[0]("foo")
		if v != `a"a"a"` {
			t.Errorf("replace modifier failed, got %v", v)
		}
	}
}

func TestParseJsonKeyWithFunc_errors(t *testing.T) {
	_, _, err := parseJsonKeyWithFunc("")
	t.Logf("err for empty key: %#v", err)
	if err == nil {
		t.Error("expected error for empty key")
	}
	_, _, err = parseJsonKeyWithFunc("foo|replace('o')")
	if err == nil {
		t.Error("expected error for invalid replace format")
	}
	_, _, err = parseJsonKeyWithFunc("foo|replace('(','a')")
	if err == nil {
		t.Error("expected error for invalid regexp")
	}
	_, _, err = parseJsonKeyWithFunc("foo|unknownmod")
	if err == nil {
		t.Error("expected error for unknown modifier")
	}
}

func TestParseJsonKey(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"foo.bar.baz", []string{"foo", "bar", "baz"}},
		{`foo."bar.baz".qux`, []string{"foo", "bar.baz", "qux"}},
		{`foo.[0].bar`, []string{"foo", "[0]", "bar"}},
		{`foo."bar.baz".[1].qux`, []string{"foo", "bar.baz", "[1]", "qux"}},
		{`"foo\".bar".baz`, []string{`foo".bar`, "baz"}},
		{`"foo\".bar"."b'az"`, []string{`foo".bar`, "b'az"}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			key, err := parseJsonKey(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(key, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, key)
			}
		})
	}
}

func TestParseJsonKeyWithFunc_cases(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
		funcs    int
	}{
		{"foo.bar.baz | tolower", []string{"foo", "bar", "baz"}, 1},
		{`foo."bar.baz".qux | tolower`, []string{"foo", "bar.baz", "qux"}, 1},
		{`foo.[0].bar | tolower`, []string{"foo", "[0]", "bar"}, 1},
		{`foo."bar.baz".[1].qux | tolower`, []string{"foo", "bar.baz", "[1]", "qux"}, 1},
		{`"foo\".bar".baz | tolower`, []string{`foo".bar`, "baz"}, 1},
		{`"foo\".bar"."b'az" | tolower`, []string{`foo".bar`, "b'az"}, 1},
		{`foo."ba|r".baz | tolower`, []string{"foo", "ba|r", "baz"}, 1},
		{`foo.'ba|r'.baz | tolower`, []string{"foo", "ba|r", "baz"}, 1},
		{`"foo".[ba|r].baz | tolower`, []string{"foo", "[ba|r]", "baz"}, 1},
		{`foo.(ba|r).baz | replace('ba|r','qux') | tolower`, []string{"foo", "(ba|r)", "baz"}, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			key, funcs, err := parseJsonKeyWithFunc(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(key, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, key)
			}
			if len(funcs) != tc.funcs {
				t.Errorf("expected %d modifier, got %d", tc.funcs, len(funcs))
			}
		})
	}
}
