package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fourls/ilo/internal/data"
	"github.com/fourls/ilo/internal/display"
	"github.com/fourls/ilo/internal/exec"
	"github.com/fourls/ilo/internal/ilofile/iloyml"
	"github.com/spf13/cobra"
)

var cmdRun = &cobra.Command{
	Use:  "run [flows...]",
	RunE: runCmdImpl,
}

var projectPath string

func init() {
	var wd, _ = os.Getwd()
	cmdRun.Flags().StringVarP(&projectPath, "project", "p", wd, "path to project definition file")
}

func runCmdImpl(cmd *cobra.Command, args []string) error {
	var stat, _ = os.Stat(projectPath)
	if stat.IsDir() {
		projectPath = filepath.Join(projectPath, "ilo.yml")
	}

	if !filepath.IsAbs(projectPath) {
		var err error
		projectPath, err = filepath.Abs(projectPath)
		if err != nil {
			return err
		}
	}

	project, err := iloyml.New(projectPath)
	if err != nil {
		return err
	}

	toolbox, _ := data.NewProdToolbox()
	log := log.New(os.Stdout, "", 0)

	executor := exec.FlowExecutor{Toolbox: *toolbox}

	observer := display.NewObserver(project, log)

	for _, flowName := range args {
		flow, exists := project.Flows[flowName]
		if !exists {
			return fmt.Errorf("no flow '%s' exists", flowName)
		}
		executor.RunFlow(flow, &observer)
	}

	return nil
}
