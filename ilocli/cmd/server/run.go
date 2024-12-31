package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"fourls.dev/ilo/ilocli/display"
	"fourls.dev/ilo/ilosrv"
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

	server := ilosrv.BuildServer()
	err := server.Run("localhost:8116")

	if err != nil {
		display.HorizontalRule{Footer: fmt.Sprintf("ERROR: %s", err.Error())}.Print(log)
	}
	return err
}
