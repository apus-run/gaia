package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/apus-run/gaia/cmd/gaia/new"
	"github.com/apus-run/gaia/cmd/gaia/rpc"
	"github.com/apus-run/gaia/cmd/gaia/run"
	"github.com/apus-run/gaia/cmd/gaia/upgrade"
)

var Cmd = &cobra.Command{
	Use:   "Gaia",
	Short: "Gaia: 基于gRPC业务开发框架",
	Long:  "Gaia: 基于gRPC业务开发框架",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func init() {
	Cmd.AddCommand(rpc.Cmd)
	Cmd.AddCommand(new.Cmd)
	Cmd.AddCommand(run.Cmd)
	Cmd.AddCommand(upgrade.Cmd)
}
func main() {
	if err := Cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
