package tool

import (
	"fmt"

	"fourls.dev/ilo/ilolib"
	"github.com/spf13/cobra"
)

var cmdToolAdd = &cobra.Command{
	Use:  "add",
	RunE: cmdToolAddImpl,
}

func cmdToolAddImpl(cmd *cobra.Command, args []string) error {
	toolbox, err := ilolib.NewProdToolbox()
	if toolbox != nil && err != nil {
		return err
	}

	for _, toolName := range args {
		info, err := toolbox.AddAuto(toolName)
		if err != nil {
			return err
		}
		fmt.Printf("Registered $%s at path '%s'\n", toolName, info.Path)
	}

	return toolbox.Save()
}
