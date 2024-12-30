package cmd

import (
	"log"
	"os"
	"path/filepath"

	"fourls.dev/ilo/ilocli/display"
	"fourls.dev/ilo/ilolib"
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

	project, err := ilolib.ReadProjectDefinition(projectPath)
	if err != nil {
		return err
	}

	toolbox, _ := ilolib.NewProdToolbox()
	log := log.New(os.Stdout, "", 0)

	executor := ilolib.ProjectExecutor{
		Definition: *project,
		Toolbox:    *toolbox,
	}

	observer := display.NewObserver(project, log)

	for _, flowName := range args {
		_, err := executor.RunFlow(flowName, &observer)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
