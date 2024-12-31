package server

import "github.com/spf13/cobra"

var CmdServer = &cobra.Command{
	Use: "server",
}

func init() {
	CmdServer.AddCommand(cmdServerRun)
}
