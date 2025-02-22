package iloyml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unicode"

	"github.com/fourls/ilo/internal/ilofile"
	"gopkg.in/yaml.v3"
)

type yamlStepDef struct {
	Echo string
	Run  string
}

type yamlProjDef struct {
	Name  string
	Flows map[string][]yamlStepDef
}

type step struct {
	text     string
	args     []string
	stepType ilofile.StepType
}

func (s step) StepType() ilofile.StepType {
	return s.stepType
}

func (s step) Args() []string {
	return s.args
}

func (s step) Message() string {
	return s.text
}

func (s step) String() string {
	return s.text
}

func New(path string) (*ilofile.Definition, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	project := ilofile.Definition{
		Path: path,
	}

	err = parseProjectDefinitionYaml(bytes, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func parseProjectDefinitionYaml(data []byte, project *ilofile.Definition) error {
	var yml yamlProjDef
	if err := yaml.Unmarshal(data, &yml); err != nil {
		return err
	}

	project.Name = yml.Name
	project.Flows = make(map[string]ilofile.Flow, len(yml.Flows))
	projectDir := filepath.Dir(project.Path)

	for flowName, flowCmds := range yml.Flows {
		flow := ilofile.Flow{
			Name:    flowName,
			Steps:   make([]ilofile.Step, len(flowCmds)),
			Project: project,
			Dir:     projectDir,
		}

		for i, line := range flowCmds {
			var stepType ilofile.StepType
			switch {
			case line.Run != "" && line.Echo == "":
				stepType = ilofile.StepRunProgram
			case line.Run == "" && line.Echo != "":
				stepType = ilofile.StepEchoMessage
			default:
				return fmt.Errorf("parse '%s' step %d: invalid type", flowName, i)
			}

			step := step{
				stepType: stepType,
			}

			switch step.stepType {
			case ilofile.StepRunProgram:
				step.text = line.Run
				parseArgsString(step.text, &step.args)
			case ilofile.StepEchoMessage:
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
