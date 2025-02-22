package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fourls/ilo/internal/data/provide"
	"github.com/fourls/ilo/internal/data/toolbox"
	"github.com/fourls/ilo/internal/display"
	"github.com/fourls/ilo/internal/server"
	"github.com/spf13/cobra"
)

var cmdServerRun = &cobra.Command{
	Use:  "run",
	RunE: cmdServerRunImpl,
}

func cmdServerRunImpl(cmd *cobra.Command, args []string) error {
	log := log.New(os.Stdout, "", 0)

	display.HorizontalRule{Header: "Ilo Automation Server"}.Print(log)

	// Print message on CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println()
		display.HorizontalRule{Footer: "Keyboard interrupt"}.Print(log)
		os.Exit(1)
	}()

	server := server.BuildServer(provide.NewConfigProvider[toolbox.Toolbox]())
	err := server.Run("localhost:8116")

	if err != nil {
		display.HorizontalRule{Footer: fmt.Sprintf("ERROR: %s", err.Error())}.Print(log)
	}
	return err
}
