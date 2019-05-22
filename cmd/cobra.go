package cmd

import (
	"etcdcc/apiserver/cmd/apiserver"
	"etcdcc/apiserver/cmd/client/file"
	"etcdcc/apiserver/cmd/client/uds"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "start"}

func init() {
	rootCmd.AddCommand(apiserver.ServeCommand)
	rootCmd.AddCommand(file.FileCommand)
	rootCmd.AddCommand(uds.UdsCommand)
}

//Cobra entrance
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
