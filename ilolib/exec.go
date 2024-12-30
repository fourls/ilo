package ilolib

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func printLines(text string, println func(string)) {
	for _, line := range strings.Split(strings.TrimRight(text, "\r\n"), "\n") {
		println(strings.TrimRight(line, "\r"))
	}
}

type ExecParams struct {
	Env       []string
	Directory string
	Observer  ExecutionObserver
	Toolbox   Toolbox
}

type StepExecutor interface {
	StepExecute(params ExecParams) error
}

type runStepExecutor struct {
	def RunFlowStep
}

func (s runStepExecutor) StepExecute(params ExecParams) error {
	args := s.def.Args()
	if len(args) < 1 {
		return errors.New("execute run step: no arguments provided")
	}

	firstArg := args[0]
	if strings.HasPrefix(firstArg, "$") {
		var info, exists = params.Toolbox.Get(firstArg[1:])
		if exists {
			firstArg = info.Path
		} else {
			return errors.New("execute run step: no tool found for substitution " + firstArg)
		}
	}

	var cmd = exec.Command(firstArg, args[1:]...)

	cmd.Env = params.Env
	cmd.Dir = params.Directory

	// todo read stderr
	var out, err = cmd.Output()

	if len(out) > 0 {
		printLines(string(out), params.Observer.StepOutput)
	}

	return err
}

type echoStepExecutor struct {
	def EchoFlowStep
}

func (s echoStepExecutor) StepExecute(params ExecParams) error {
	printLines(s.def.Message(), params.Observer.StepOutput)
	return nil
}

func BuildDefaultExecutor(step FlowStep) StepExecutor {
	switch step.StepType() {
	case StepEchoMessage:
		return echoStepExecutor{step.(EchoFlowStep)}
	case StepRunProgram:
		return runStepExecutor{step.(RunFlowStep)}
	default:
		return nil
	}
}

type ProjectExecutor struct {
	Definition ProjectDefinition
	Toolbox    Toolbox
}

type FlowExecutionError struct {
	FlowName string
	Message  string
}

func (e FlowExecutionError) Error() string {
	return fmt.Sprintf("execute flow '%s': %s", e.FlowName, e.Message)
}

type ExecutionObserver interface {
	FlowEntered(f *FlowDef)
	FlowPassed()
	FlowFailed()

	StepEntered(s *FlowStep)
	StepOutput(text string)
	StepPassed()
	StepFailed(err error)
}

func (e ProjectExecutor) runStep(flow FlowDef, index int, params ExecParams, buildExecutor func(FlowStep) StepExecutor) error {
	executor := buildExecutor(flow.Steps[index])
	if executor == nil {
		// Unknown step type
		return FlowExecutionError{
			FlowName: flow.Name,
			Message:  fmt.Sprintf("step %d is unknown and cannot be processed", index),
		}
	}

	err := executor.StepExecute(params)
	if err != nil {
		return FlowExecutionError{
			FlowName: flow.Name,
			Message:  fmt.Sprintf("step %d failed: %v", index, err),
		}
	}

	return nil
}

func (e ProjectExecutor) RunFlow(
	name string,
	observer ExecutionObserver,
	buildExecutor func(FlowStep) StepExecutor,
) (bool, error) {
	flow, ok := e.Definition.Flows[name]
	if !ok {
		return false, FlowExecutionError{
			name,
			"flow does not exist",
		}
	}

	observer.FlowEntered(&flow)

	baseParams := ExecParams{
		Env:       os.Environ(),
		Directory: e.Definition.Dir,
		Observer:  observer,
		Toolbox:   e.Toolbox,
	}

	success := true

	for i := range flow.Steps {
		observer.StepEntered(&flow.Steps[i])
		err := e.runStep(flow, i, baseParams, buildExecutor)

		if err != nil {
			observer.StepFailed(err)
			success = false
			break
		} else {
			observer.StepPassed()
		}
	}

	if success {
		observer.FlowPassed()
	} else {
		observer.FlowFailed()
	}

	return success, nil
}
