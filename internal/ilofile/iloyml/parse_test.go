package iloyml

import (
	"reflect"
	"testing"

	"github.com/fourls/ilo/internal/ilofile"
)

func TestParseFlows(t *testing.T) {
	var data = []byte(`
flows:
  foo:
    - echo: Starting foo
    - run: cmd -abc "this is a foo text"
    - echo: Finishing foo
  bar:
    - echo: Doing "bar" now`)

	var def ilofile.Definition
	if err := parseProjectDefinitionYaml(data, &def); err != nil {
		t.Fatalf("got: %v, want: nil", err)
	}

	if len(def.Flows) != 2 {
		t.Fatalf("got: %d flows, want: 2 flows", len(def.Flows))
	}

	var flow = def.Flows["foo"]
	if flow.Name != "foo" || len(flow.Steps) != 3 {
		t.Fatalf("got: %d steps in flow, want: 3 steps in flow", len(flow.Steps))
	}

	var expectedStep = step{stepType: ilofile.StepEchoMessage, text: "Starting foo"}

	if !reflect.DeepEqual(flow.Steps[0], expectedStep) {
		t.Fatalf("got: %s, want: %s", flow.Steps[0], expectedStep)
	}

	expectedStep = step{
		stepType: ilofile.StepRunProgram,
		text:     "cmd -abc \"this is a foo text\"",
		args:     []string{"cmd", "-abc", "this is a foo text"},
	}

	if !reflect.DeepEqual(flow.Steps[1], expectedStep) {
		t.Fatalf("got: %s, want: %s", flow.Steps[1], expectedStep)
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
