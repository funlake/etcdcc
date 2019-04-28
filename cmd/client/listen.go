package client

import (
	"etcdcc/apiserver/pkg/client"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/cobra"
)

var (
	mod            string
	EtcdCertFile   string
	EtcdKeyFile    string
	EtcdCaFile     string
	EtcdServerName string
	EtcdHosts      string
)

func init() {
	var cp = ClientCommand.PersistentFlags()
	cp.StringVar(&mod, "gmod", "global", "Name of prefix of global module")
	cp.StringVar(&mod, "mod", "global", "Name of prefix of current module")
	cp.StringVar(&EtcdCertFile, "c", "/keys/ca.pem", "Cert file for etcd connection")
	cp.StringVar(&EtcdKeyFile, "k", "/keys/ca-key.pem", "Key file for etcd connection")
	cp.StringVar(&EtcdCaFile, "ca", "/keys/ca.crt", "Ca file for etcd connection")
	cp.StringVar(&EtcdServerName, "sn", "etcchebao", "Hostname for ssl verification")
	cp.StringVar(&EtcdHosts, "hosts", "127.0.0.1:2379", "Hosts of etcd server")
	if cobra.MarkFlagRequired(cp, "mod") != nil ||
		cobra.MarkFlagRequired(cp, "hosts") != nil {
		//cobra.MarkFlagRequired(sp, "k") != nil ||
		//cobra.MarkFlagRequired(sp, "ca") != nil {
		log.Error("Fail to set required")
	}
}

var ClientCommand = &cobra.Command{
	Use:   "client",
	Short: "Listining config changes & modified local configuration",
	Run: func(cmd *cobra.Command, args []string) {
		etcd.Adapter.Connect(etcd.Adapter{},EtcdHosts,EtcdCertFile,EtcdKeyFile,EtcdCaFile,EtcdServerName)
		log.Info("Successfully connected to etcd server")
		wc := &client.EtcdClientWatcher{}
		wc.KeepEyesOnKeyWithPrefix(fmt.Sprintf("dev/%s", mod), clientv3.WithPrefix())
	},
}
