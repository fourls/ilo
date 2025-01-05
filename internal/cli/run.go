package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fourls/ilo/internal/data/provide"
	"github.com/fourls/ilo/internal/data/toolbox"
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

	provider := provide.NewConfigProvider[toolbox.Toolbox]()
	toolbox, _ := provider.Load(
		"toolbox",
		provide.YamlUnmarshal[toolbox.Toolbox])

	log := log.New(os.Stdout, "", 0)
	observer := display.NewObserver(project, log)

	for _, flowName := range args {
		flow, exists := project.Flows[flowName]
		if !exists {
			return fmt.Errorf("no flow '%s' exists", flowName)
		}
		exec.RunFlow(flow, exec.RunStep, *toolbox, &observer)
	}

	return nil
}
