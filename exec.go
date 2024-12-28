package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunFlow(flow FlowDef) error {
	var basePrefix = fmt.Sprintf("[%s] ", flow.Name)

	var flowLog = log.New(os.Stdout, basePrefix, 0)
	flowLog.Println("START")

	for i, step := range flow.Steps {
		flowLog.SetPrefix(fmt.Sprintf("%s%d: ", basePrefix, i))
		var cmd = exec.Command(step.Args[0], step.Args[1:]...)
		var out, err = cmd.Output()

		if len(out) > 0 {
			for _, line := range strings.Split(strings.TrimRight(string(out), "\r\n"), "\n") {
				flowLog.Println(strings.TrimRight(line, "\r"))
			}
		}

		if err != nil {
			flowLog.Printf("(error) %s\n", err.Error())
			flowLog.Println("FAIL")
			return errors.New("step failed")
		}
	}

	flowLog.SetPrefix(basePrefix)
	flowLog.Println("PASS")
	return nil
}
