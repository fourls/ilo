package server

import (
	"fourls.dev/ilo/ilosrv"
	"github.com/spf13/cobra"
)

var cmdServerRun = &cobra.Command{
	Use:  "run",
	RunE: cmdServerRunImpl,
}

func cmdServerRunImpl(cmd *cobra.Command, args []string) error {
	server := ilosrv.BuildServer()
	return server.Run("localhost:8116")
}
