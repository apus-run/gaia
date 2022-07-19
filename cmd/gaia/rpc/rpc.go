package rpc

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "rpc",
	Short: "Rpc project",
	Long:  "Rpc project. Example: gaia rpc",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {

}
