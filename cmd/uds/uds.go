/**
unix domain socket for easily integration
 */
package uds

import (
	"etcdcc/apiserver/pkg/client"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/spf13/cobra"
)
var (
	mod            string
	gmod           string
	etcdCertFile   string
	etcdKeyFile    string
	etcdCaFile     string
	etcdServerName string
	etcdHosts      string
	retrySeconds   int
	storeDir       string
)
func init(){
	var cp = UdsCommand.PersistentFlags()
	cp.StringVar(&gmod, "gmod", "global", "Name of prefix of global module")
	cp.StringVar(&mod, "mod", "global", "Name of prefix of current module")
	cp.StringVar(&etcdCertFile, "c", "/keys/ca.pem", "Cert file for etcd connection")
	cp.StringVar(&etcdKeyFile, "k", "/keys/ca-key.pem", "Key file for etcd connection")
	cp.StringVar(&etcdCaFile, "ca", "/keys/ca.crt", "Ca file for etcd connection")
	cp.StringVar(&etcdServerName, "sn", "etcchebao", "Hostname for ssl verification")
	cp.StringVar(&etcdHosts, "hosts", "127.0.0.1:2379", "Hosts of etcd server")
	cp.IntVar(&retrySeconds,"retrySeconds",3,"Fails retry in ? seconds")
	cp.StringVar(&storeDir,"storeDir","/tmp","Directory of config file")

	if cobra.MarkFlagRequired(cp, "mod") != nil ||
		cobra.MarkFlagRequired(cp, "hosts") != nil {
		//cobra.MarkFlagRequired(sp, "k") != nil ||
		//cobra.MarkFlagRequired(sp, "ca") != nil {
		log.Error("Fail to set required")
	}
}
var UdsCommand = &cobra.Command{
	Use:   "uds",
	Short: "Listening config changes & server on uds",
	Run: func(cmd *cobra.Command, args []string) {
		etcd.Adapter.Connect(etcd.Adapter{}, etcdHosts, etcdCertFile, etcdKeyFile, etcdCaFile, etcdServerName)
		log.Info("Successfully connected to etcd server[uds]")
		wc := &client.EtcdUdsWatcher{}
		go wc.Serve()
		wc.KeepEyesOnKeyWithPrefix(fmt.Sprintf("dev/%s", mod))

	},
}
