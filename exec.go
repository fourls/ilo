package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunExecStep(program string, args []string, logger *log.Logger) bool {
	var cmd = exec.Command(program, args...)
	var out, err = cmd.Output()

	if len(out) > 0 {
		for _, line := range strings.Split(strings.TrimRight(string(out), "\r\n"), "\n") {
			logger.Println(strings.TrimRight(line, "\r"))
		}
	}

	if err != nil {
		logger.Printf("(error) %s\n", err.Error())
		return false
	}
	return true
}

func RunFlow(flow FlowDef) error {
	var basePrefix = fmt.Sprintf("[%s] ", flow.Name)

	var logger = log.New(os.Stdout, basePrefix, 0)
	logger.Println("BEGIN")

	for i, step := range flow.Steps {
		logger.SetPrefix(fmt.Sprintf("%s%d: ", basePrefix, i))

		if step.Echo != "" {
			logger.Println(step.Echo)
		} else if len(step.Run) > 0 {
			if !RunExecStep(step.Run[0], step.Run[1:], logger) {
				logger.Println("END (FAIL)")
				return errors.New("step failed")
			}
		}
	}

	logger.SetPrefix(basePrefix)
	logger.Println("END (PASS)")
	return nil
}
