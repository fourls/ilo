package cmd

import (
	"fmt"
	"os"

	"fourls.dev/ilo/ilocli/cmd/tool"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "ilo",
	Short: "A simple task runner.",
}

func init() {
	cmdRoot.AddCommand(cmdRun)
	cmdRoot.AddCommand(tool.CmdTool)
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
