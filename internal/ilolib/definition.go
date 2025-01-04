package ilolib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unicode"

	"gopkg.in/yaml.v3"
)

type FlowStep interface {
	StepType() StepType
	String() string
}

type RunFlowStep interface {
	Args() []string
}

type EchoFlowStep interface {
	Message() string
}

type StepType int

const (
	StepRunProgram StepType = iota
	StepEchoMessage
)

type flowStep struct {
	text     string
	args     []string
	stepType StepType
}

func (s flowStep) StepType() StepType {
	return s.stepType
}

func (s flowStep) Args() []string {
	return s.args
}

func (s flowStep) Message() string {
	return s.text
}

func (s flowStep) String() string {
	return s.text
}

type Flow struct {
	Name    string
	Dir     string
	Steps   []FlowStep
	Project *ProjectDefinition
}

type ProjectDefinition struct {
	Name  string
	Path  string
	Flows map[string]Flow
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
	project.Flows = make(map[string]Flow, len(yml.Flows))
	projectDir := filepath.Dir(project.Path)

	for flowName, flowCmds := range yml.Flows {
		flow := Flow{
			Name:    flowName,
			Steps:   make([]FlowStep, len(flowCmds)),
			Project: project,
			Dir:     projectDir,
		}

		for i, line := range flowCmds {
			var stepType StepType
			switch {
			case line.Run != "" && line.Echo == "":
				stepType = StepRunProgram
			case line.Run == "" && line.Echo != "":
				stepType = StepEchoMessage
			default:
				return fmt.Errorf("parse '%s' step %d: invalid type", flowName, i)
			}

			step := flowStep{
				stepType: stepType,
			}

			switch step.stepType {
			case StepRunProgram:
				step.text = line.Run
				parseArgsString(step.text, &step.args)
			case StepEchoMessage:
				step.text = line.Echo
			}

			flow.Steps[i] = step
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
