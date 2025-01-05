package exec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fourls/ilo/internal/data/toolbox"
	"github.com/fourls/ilo/internal/ilofile"
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
	Toolbox   toolbox.Toolbox
}

func doRunStep(step ilofile.RunFlowStep, params ExecParams) error {
	args := step.Args()
	if len(args) < 1 {
		return errors.New("execute run step: no arguments provided")
	}

	firstArg := args[0]
	if strings.HasPrefix(firstArg, "$") {
		var path, exists = params.Toolbox[firstArg[1:]]
		if exists {
			firstArg = path
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

func RunStep(step ilofile.Step, params ExecParams) error {
	switch step.StepType() {
	case ilofile.StepEchoMessage:
		printLines(step.(ilofile.EchoFlowStep).Message(), params.Observer.StepOutput)
		return nil
	case ilofile.StepRunProgram:
		return doRunStep(step.(ilofile.RunFlowStep), params)
	default:
		return errors.New("step failed: Unknown step type")
	}
}

type StepExecutorFunc func(ilofile.Step, ExecParams) error

type FlowExecutionError struct {
	FlowName string
	Message  string
}

func (e FlowExecutionError) Error() string {
	return fmt.Sprintf("execute flow '%s': %s", e.FlowName, e.Message)
}

type ExecutionObserver interface {
	FlowEntered(f *ilofile.Flow)
	FlowPassed()
	FlowFailed()

	StepEntered(s ilofile.Step)
	StepOutput(text string)
	StepPassed()
	StepFailed(err error)
}

type noOpObserver struct{}

func (o noOpObserver) FlowEntered(f *ilofile.Flow) {}
func (o noOpObserver) FlowPassed()                 {}
func (o noOpObserver) FlowFailed()                 {}
func (o noOpObserver) StepEntered(s ilofile.Step)  {}
func (o noOpObserver) StepOutput(text string)      {}
func (o noOpObserver) StepPassed()                 {}
func (o noOpObserver) StepFailed(err error)        {}

// RunFlow executes all steps in the specified flow using stepExecutor,
// and reports whether all steps were executed successfully.
// If stepExecutor is omitted, steps will not be executed.
// If provided, the observer will be called alongside various milestones,
// see ExecutionObserver for more information.
func RunFlow(
	flow ilofile.Flow,
	stepExecutor StepExecutorFunc,
	toolbox toolbox.Toolbox,
	observer ExecutionObserver,
) bool {
	if stepExecutor == nil {
		stepExecutor = func(ilofile.Step, ExecParams) error { return nil }
	}
	if observer == nil {
		observer = noOpObserver{}
	}

	observer.FlowEntered(&flow)

	baseParams := ExecParams{
		Env:       os.Environ(),
		Directory: flow.Dir,
		Observer:  observer,
		Toolbox:   toolbox,
	}

	success := true

	for i := range flow.Steps {
		observer.StepEntered(flow.Steps[i])

		err := stepExecutor(flow.Steps[i], baseParams)

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

	return success
}
