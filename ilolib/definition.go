package ilolib

import (
	"errors"
	"os"
	"path/filepath"
	"unicode"

	"gopkg.in/yaml.v3"
)

type FlowStepDef struct {
	Run  []string
	Echo string
}

type FlowDef struct {
	Name  string
	Steps []FlowStepDef
}

type ProjectDefinition struct {
	Name  string
	Path  string
	Dir   string
	Flows map[string]FlowDef
}

type YamlStepDef struct {
	Echo string
	Run  string
}

type YamlProjDef struct {
	Name  string
	Flows map[string][]YamlStepDef
}

func ReadProjectDefinition(path string) (*ProjectDefinition, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	project := ProjectDefinition{
		Path: path,
		Dir:  filepath.Dir(path),
	}

	err = parseProjectDefinitionYaml(bytes, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func parseProjectDefinitionYaml(data []byte, project *ProjectDefinition) error {
	var yml YamlProjDef
	if err := yaml.Unmarshal(data, &yml); err != nil {
		return err
	}

	project.Name = yml.Name
	project.Flows = make(map[string]FlowDef, len(yml.Flows))

	for flowName, flowCmds := range yml.Flows {
		var flow FlowDef

		flow.Name = flowName
		flow.Steps = make([]FlowStepDef, len(flowCmds))

		for i, line := range flowCmds {
			if line.Run != "" && line.Echo != "" {
				return errors.New("flow step contains both run: and echo: commands")
			}

			switch {
			case line.Run != "":
				parseArgsString(line.Run, &flow.Steps[i].Run)
			case line.Echo != "":
				flow.Steps[i].Echo = line.Echo
			}
		}

		project.Flows[flowName] = flow
	}

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
