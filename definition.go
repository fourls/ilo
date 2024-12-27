package main

import (
	"encoding/json"
	"errors"
	"strings"
	"unicode"
)

type ToolRef struct {
	Name    string
	Version string
}

type FlowStep struct {
	Args []string
}

type FlowDef struct {
	Name  string
	Steps []FlowStep
}

type ProjDef struct {
	Tools []ToolRef
	Flows []FlowDef
}

type JsonProjDef struct {
	Tools []string
	Flows map[string][]string
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

func ParseProjDef(bytes []byte) (*ProjDef, error) {
	var data JsonProjDef
	var err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	// Parse tools
	tools := make([]ToolRef, len(data.Tools))
	for i, str := range data.Tools {
		parseToolRef(str, &tools[i])
	}

	// Parse flows
	var flows = make([]FlowDef, len(data.Flows))
	var flowIndex = 0
	for key, value := range data.Flows {
		var flow = &flows[flowIndex]
		flow.Name = key

		for i, line := range value {
			parseArgsString(line, &flow.Steps[i].Args)
		}
		flowIndex++
	}

	return &ProjDef{
		Tools: tools,
	}, nil
}
