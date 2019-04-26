package cmd

import (
	"etcdcc/apiserver/cmd/apiserver"
	"etcdcc/apiserver/cmd/client"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)
var RootCmd = &cobra.Command{Use: "start"}
func init(){
	RootCmd.AddCommand(apiserver.ServeCommand)
	RootCmd.AddCommand(client.ClientCommand)
}
func Execute()  {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
