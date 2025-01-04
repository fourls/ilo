package tool

import (
	"github.com/spf13/cobra"
)

var CmdTool = &cobra.Command{
	Use: "tool",
}

func init() {
	CmdTool.AddCommand(cmdToolAdd)
}
