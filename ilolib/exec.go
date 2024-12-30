package ilolib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func printLines(text string, logger *log.Logger) {
	for _, line := range strings.Split(strings.TrimRight(text, "\r\n"), "\n") {
		logger.Println(strings.TrimRight(line, "\r"))
	}
}

type execParams struct {
	Env       []string
	Directory string
	Logger    *log.Logger
	Toolbox   Toolbox
}

type executor interface {
	execute(params execParams) error
}

type runStepExecutor struct {
	def RunFlowStep
}

func (s runStepExecutor) execute(params execParams) error {
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

	var out, err = cmd.Output()

	if len(out) > 0 {
		printLines(string(out), params.Logger)
	}

	return err
}

type echoStepExecutor struct {
	def EchoFlowStep
}

func (s echoStepExecutor) execute(params execParams) error {
	printLines(s.def.Message(), params.Logger)
	return nil
}

func buildExecutor(step FlowStep) executor {
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

func (e ProjectExecutor) runStep(flow FlowDef, index int, params execParams) error {
	var err error
	defer func() {
		if err != nil {
			params.Logger.Println("ERROR: " + err.Error())
		}
	}()

	executor := buildExecutor(flow.Steps[index])
	if executor == nil {
		// Unknown step type
		return FlowExecutionError{
			FlowName: flow.Name,
			Message:  fmt.Sprintf("step %d is unknown and cannot be processed", index),
		}
	}

	err = executor.execute(params)
	if err != nil {
		return FlowExecutionError{
			FlowName: flow.Name,
			Message:  fmt.Sprintf("step %d failed: %v", index, err),
		}
	}

	return nil
}

func (e ProjectExecutor) RunFlow(name string, log *log.Logger) (bool, error) {
	flow, ok := e.Definition.Flows[name]
	if !ok {
		return false, FlowExecutionError{
			name,
			"flow does not exist",
		}
	}

	flowId := fmt.Sprintf("%s / %s", e.Definition.Name, flow.Name)

	HorizontalRule{Header: flowId}.Print(log)

	timeStarted := time.Now()

	baseParams := execParams{
		Env:       os.Environ(),
		Directory: e.Definition.Dir,
		Logger:    log,
		Toolbox:   e.Toolbox,
	}

	success := true

	for i := range flow.Steps {
		err := e.runStep(flow, i, baseParams)

		if err != nil {
			success = false
			break
		}
	}

	duration := time.Since(timeStarted).Round(time.Millisecond)

	var status string
	if success {
		status = fmt.Sprintf("PASSED in %s", duration)
	} else {
		status = fmt.Sprintf("FAILED after %s", duration)
	}

	HorizontalRule{Footer: status}.Print(log)

	return success, nil
}
