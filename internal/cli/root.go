package cli

import (
	"fmt"
	"os"

	"github.com/fourls/ilo/internal/cli/server"
	"github.com/fourls/ilo/internal/cli/tool"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "ilo",
	Short: "A simple task runner.",
}

func init() {
	cmdRoot.AddCommand(cmdRun)
	cmdRoot.AddCommand(tool.CmdTool)
	cmdRoot.AddCommand(server.CmdServer)
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
