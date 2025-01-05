package tool

import (
	"fmt"

	"github.com/fourls/ilo/internal/data/provide"
	"github.com/fourls/ilo/internal/data/toolbox"
	"github.com/spf13/cobra"
)

var cmdToolAdd = &cobra.Command{
	Use:  "add",
	RunE: cmdToolAddImpl,
}

func cmdToolAddImpl(cmd *cobra.Command, args []string) error {
	provider := provide.NewConfigProvider[toolbox.Toolbox]()
	toolbox, err := provider.Load("toolbox",
		provide.YamlUnmarshal[toolbox.Toolbox])
	if err != nil {
		return err
	}

	if *toolbox == nil {
		// We want to be able to update the toolbox
		*toolbox = make(map[string]string)
	}

	for _, name := range args {
		err := toolbox.FindAndAdd(name)
		if err != nil {
			return err
		}
		fmt.Printf("Registered $%s at path '%s'\n", name, (*toolbox)[name])
	}

	return provider.Save("toolbox", toolbox, provide.YamlMarshal)
}
