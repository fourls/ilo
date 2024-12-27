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
		InNone = iota
		InUnescaped
		InSingleQuote
		InDoubleQuote
	)

	var inState = InNone
	var start = 0

	var setStart = func(index int, state int) {
		start = index
		inState = state
	}

	var setEnd = func(index int) {
		inState = InNone
		*out = append(*out, line[start:index])
	}

	for i, char := range line {
		switch {
		case char == '\'' && inState == InNone:
			setStart(i+1, InSingleQuote)
		case char == '"' && inState == InNone:
			setStart(i+1, InDoubleQuote)
		case char == '\'' && inState == InSingleQuote:
			setEnd(i)
		case char == '"' && inState == InDoubleQuote:
			setEnd(i)
		case unicode.IsSpace(char) && inState == InUnescaped:
			setEnd(i)
		case !unicode.IsSpace(char) && inState == InNone:
			setStart(i, InUnescaped)
		}
	}

	switch inState {
	case InUnescaped:
		setEnd(len(line))
		return nil
	case InNone:
		return nil
	default:
		return errors.New("unterminated string literal in flow step")
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
