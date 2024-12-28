package main

import (
	"errors"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type ToolRef struct {
	Name    string
	Version string
}

type FlowStep struct {
	Run  []string
	Echo string
}

type FlowDef struct {
	Name  string
	Steps []FlowStep
}

type ProjDef struct {
	Flows []FlowDef
}

type JsonProjDef struct {
	Tools []string
	Flows map[string][]string
}

type YamlStepDef struct {
	Echo string
	Run  string
}

type YamlProjDef struct {
	Flows map[string][]YamlStepDef
}

func parseYamlProjDef(data []byte) (*ProjDef, error) {
	var yml YamlProjDef
	var err = yaml.Unmarshal(data, &yml)
	if err != nil {
		return nil, err
	}

	var flows = make([]FlowDef, len(yml.Flows))
	var flowIndex = 0
	for flowName, flowCmds := range yml.Flows {
		var flow = &flows[flowIndex]

		flow.Name = flowName
		flow.Steps = make([]FlowStep, len(flowCmds))

		for i, line := range flowCmds {
			if line.Run != "" && line.Echo != "" {
				return nil, errors.New("flow step contains both run: and echo: commands")
			}

			switch {
			case line.Run != "":
				parseArgsString(line.Run, &flow.Steps[i].Run)
			case line.Echo != "":
				flow.Steps[i].Echo = line.Echo
			}
		}
		flowIndex++
	}

	return &ProjDef{Flows: flows}, nil
}

func parseToolRef(value string, out *ToolRef) error {
	split := strings.Split(value, "@")

	if len(split) == 0 || len(split) > 2 {
		return errors.New("Invalid tool reference " + value)
	}

	toolName := split[0]
	var version string

	if len(split) == 2 {
		version = split[1]
	}

	*out = ToolRef{Name: toolName, Version: version}
	return nil
}

func parseArgsString(line string, out *[]string) error {
	const (
		None = iota
		Literal
		SingleQuoted
		DoubleQuoted
	)

	// We want to create the args array regardless of whether there are actually any args
	*out = make([]string, 0)

	var argType = None
	var start = 0

	var setStart = func(index int, state int) {
		start = index
		argType = state
	}

	var setEnd = func(index int) {
		argType = None
		*out = append(*out, line[start:index])
	}

	for i, char := range line {
		if argType == None {
			switch {
			case char == '\'':
				setStart(i+1, SingleQuoted)
			case char == '"':
				setStart(i+1, DoubleQuoted)
			case !unicode.IsSpace(char):
				setStart(i, Literal)
			}
		} else if argType == Literal {
			switch {
			case char == '\'':
				setEnd(i)
				setStart(i+1, SingleQuoted)
			case char == '"':
				setEnd(i)
				setStart(i+1, DoubleQuoted)
			case unicode.IsSpace(char):
				setEnd(i)
			}
		} else if argType == SingleQuoted && char == '\'' {
			setEnd(i)
		} else if argType == DoubleQuoted && char == '"' {
			setEnd(i)
		}
	}

	if argType == Literal {
		setEnd(len(line))
	}

	if argType != None {
		return errors.New("unterminated string literal")
	} else {
		return nil
	}
}
