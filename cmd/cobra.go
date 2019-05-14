package cmd

import (
	"etcdcc/apiserver/cmd/apiserver"
	"etcdcc/apiserver/cmd/client/file"
	"etcdcc/apiserver/cmd/client/uds"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)
var RootCmd = &cobra.Command{Use: "start"}
func init(){
	RootCmd.AddCommand(apiserver.ServeCommand)
	RootCmd.AddCommand(file.FileCommand)
	RootCmd.AddCommand(uds.UdsCommand)
}
func Execute()  {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
