package ilolib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	Args []string
}

func (s runStepExecutor) execute(params execParams) error {
	var firstArg = s.Args[0]
	if strings.HasPrefix(firstArg, "$") {
		var info, exists = params.Toolbox.Tools[firstArg[1:]]
		if exists {
			firstArg = info.Path
		} else {
			return errors.New("no tool found for substitution " + firstArg)
		}
	}

	var cmd = exec.Command(firstArg, s.Args[1:]...)

	cmd.Env = params.Env
	cmd.Dir = params.Directory

	var out, err = cmd.Output()

	if len(out) > 0 {
		printLines(string(out), params.Logger)
	}

	return err
}

type echoStepExecutor struct {
	Message string
}

func (s echoStepExecutor) execute(params execParams) error {
	printLines(s.Message, params.Logger)
	return nil
}

func buildExecutor(step FlowStepDef) (executor, error) {
	if step.Echo != "" {
		return echoStepExecutor{step.Echo}, nil
	} else if len(step.Run) > 0 {
		return runStepExecutor{step.Run}, nil
	} else {
		return nil, errors.New("invalid step")
	}
}

func logFlowFinish(logger *log.Logger, prefix string, status string) {
	logger.SetPrefix(prefix)
	logger.Println(status)
}

func ExecuteFlow(flow FlowDef, toolbox Toolbox) error {
	var basePrefix = fmt.Sprintf("%s.", flow.Name)

	var logger = log.New(os.Stdout, basePrefix, 0)
	logger.Println("BEGIN")

	var baseEnv = os.Environ()
	var defaultDir, err = os.Getwd()
	if err != nil {
		return err
	}

	for i, step := range flow.Steps {
		logger.SetPrefix(fmt.Sprintf("%s%d: ", basePrefix, i))

		var executor, err = buildExecutor(step)
		if err != nil {
			logFlowFinish(logger, basePrefix, fmt.Sprintf("ERROR: parsing step %d: %s", i, err))
			return err
		}

		err = executor.execute(execParams{
			Env:       baseEnv,
			Directory: defaultDir,
			Logger:    logger,
			Toolbox:   toolbox,
		})

		if err != nil {
			logFlowFinish(logger, basePrefix, fmt.Sprintf("FAIL at step %d: %s", i, err))
			return err
		}
	}

	logFlowFinish(logger, basePrefix, "PASS")
	return nil
}
