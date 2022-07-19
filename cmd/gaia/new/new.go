package new

import (
	"github.com/spf13/cobra"

	"github.com/apus-run/gaia/cmd/gaia/new/model"
	"github.com/apus-run/gaia/cmd/gaia/new/service"
)

// Cmd represents the new command.
var Cmd = &cobra.Command{
	Use:   "new",
	Short: "Generate the new files",
	Long:  "Generate the new files.",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func init() {
	Cmd.AddCommand(model.Cmd)
	Cmd.AddCommand(service.Cmd)
}
