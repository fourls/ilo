package cmd

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

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

func printHeader(name string, path string) {
	var log = log.New(os.Stdout, "", 0)
	log.Printf(`╔%s╗`, strings.Repeat("═", 50))
	log.Printf("║ ▒▒ %-45s ║\n", name)
	log.Printf(`╟%s╢`, strings.Repeat("─", 50))
	log.Printf("║ %-48s ║\n", path)
	log.Printf(`╚%s╝`, strings.Repeat("═", 50))
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

	file_contents, err := os.ReadFile(projectPath)
	if err != nil {
		return err
	}

	project, err := ilolib.ParseYamlProjDef(file_contents)
	if err != nil {
		return err
	}

	printHeader(project.Name, projectPath)

	toolbox, _ := ilolib.NewProdToolbox()

	for _, flow := range project.Flows {
		if len(args) == 0 || slices.Contains(args, flow.Name) {
			if err = ilolib.ExecuteFlow(flow, *toolbox); err != nil {
				return err
			}
		}
	}

	return nil
}
