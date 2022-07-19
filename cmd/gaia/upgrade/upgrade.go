package upgrade

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the gaia tools",
	Long:  "Upgrade the gaia tools. Example: gaia upgrade",
	Run:   Run,
}

// Run upgrade the gaia tools.
func Run(cmd *cobra.Command, args []string) {

}
