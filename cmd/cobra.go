package cmd

import (
	"fmt"
	"github.com/funlake/etcdcc/cmd/client/file"
	"github.com/funlake/etcdcc/cmd/client/uds"
	"github.com/funlake/etcdcc/cmd/server"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "start"}

func init() {
	rootCmd.AddCommand(server.ServeCommand)
	rootCmd.AddCommand(file.FileCommand)
	rootCmd.AddCommand(uds.UdsCommand)
}

//Execute : Cobra entrance
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
