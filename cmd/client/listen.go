package client

import (
	"etcdcc/apiserver/pkg/log"
	"github.com/spf13/cobra"
)

var ClientCommand = &cobra.Command{
	Use: "listen",
	Short: "Listining config changes & modified local configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("This is listening sdk")
	},
}