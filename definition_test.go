package main

import (
	"reflect"
	"testing"
)

func TestParseToolRef(t *testing.T) {
	var tests = []struct {
		input string
		want  ToolRef
	}{
		{"foo", ToolRef{Name: "foo", Version: ""}},
		{"foo@bar", ToolRef{Name: "foo", Version: "bar"}},
		{"foo bar@baz flarp", ToolRef{Name: "foo bar", Version: "baz flarp"}},
		{"@bar", ToolRef{Name: "", Version: "bar"}},
	}

	for _, tc := range tests {
		var toolRef ToolRef
		var err = parseToolRef(tc.input, &toolRef)
		if err != nil || toolRef.Name != tc.want.Name || toolRef.Version != tc.want.Version {
			t.Fatalf(
				"got: %s, %v, want: %s, nil",
				toolRef,
				err,
				tc.want,
			)
		}
	}
}

func TestParseArgsString(t *testing.T) {
	var tests = []struct {
		input    string
		expected []string
	}{
		{
			`foo bar 'baz "bonk" flarp' "blinky's" bonk`,
			[]string{
				"foo",
				"bar",
				"baz \"bonk\" flarp",
				"blinky's",
				"bonk",
			},
		},
		{
			`foo's bar'`,
			[]string{
				"foo",
				"s bar",
			},
		},
		{
			``,
			[]string{},
		},
	}

	for _, tc := range tests {
		var args []string
		var err = parseArgsString(tc.input, &args)

		if err != nil || (!reflect.DeepEqual(args, tc.expected)) {
			t.Fatalf(
				"got: %s, %v, want: %s, nil",
				args,
				err,
				tc.expected,
			)
		}
	}

}
