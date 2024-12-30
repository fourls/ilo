package tool

import (
	"fmt"

	"fourls.dev/ilo/ilolib"
	"github.com/hairyhenderson/go-which"
	"github.com/spf13/cobra"
)

var cmdToolAdd = &cobra.Command{
	Use:  "add",
	RunE: cmdToolAddImpl,
}

func cmdToolAddImpl(cmd *cobra.Command, args []string) error {
	toolbox, err := ilolib.GetToolbox()
	if err != nil {
		return err
	}

	for _, toolName := range args {
		path := which.Which(toolName)
		if path == "" {
			return fmt.Errorf("add tool '%s': could not find on PATH", toolName)
		}

		info := ilolib.ToolInfo{
			Name: toolName,
			Path: path,
		}

		toolbox.Tools[toolName] = info
		fmt.Printf("Registered $%s at path '%s'\n", toolName, info.Path)
	}

	return ilolib.UpdateToolbox(toolbox)
}
