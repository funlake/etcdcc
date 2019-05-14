/**
unix domain socket for easily integration
 */
package uds

import (
	"etcdcc/apiserver/pkg/client"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"etcdcc/apiserver/pkg/log"
	"github.com/spf13/cobra"
)
var (
	prefix         string
	etcdCertFile   string
	etcdKeyFile    string
	etcdCaFile     string
	etcdServerName string
	etcdHosts      string
	retrySeconds   int
	sockFile       string
)
func init(){
	var cp = UdsCommand.PersistentFlags()
	cp.StringVar(&prefix, "prefix", "global", "Prefix of configuration in etcd")
	cp.StringVar(&etcdCertFile, "c", "/keys/ca.pem", "Cert file for etcd connection")
	cp.StringVar(&etcdKeyFile, "k", "/keys/ca-key.pem", "Key file for etcd connection")
	cp.StringVar(&etcdCaFile, "ca", "/keys/ca.crt", "Ca file for etcd connection")
	cp.StringVar(&etcdServerName, "sn", "", "ServerName for ssl verification")
	cp.StringVar(&etcdHosts, "hosts", "127.0.0.1:2379", "Hosts of etcd server")
	cp.StringVar(&sockFile,"sock","/dev/shm/etcdcc.sock","Unix domain socket file")

	if cobra.MarkFlagRequired(cp, "prefix") != nil ||
		cobra.MarkFlagRequired(cp, "hosts") != nil {
		//cobra.MarkFlagRequired(sp, "k") != nil ||
		//cobra.MarkFlagRequired(sp, "ca") != nil {
		log.Error("Fail to set required")
	}
}
var UdsCommand = &cobra.Command{
	Use:   "uds",
	Short: "Listening config changes & server on unix domain socket",
	Run: func(cmd *cobra.Command, args []string) {
		etcd.Adapter.Connect(etcd.Adapter{}, etcdHosts, etcdCertFile, etcdKeyFile, etcdCaFile, etcdServerName)
		log.Info("Successfully connected to etcd server[uds]")
		wc := &client.EtcdUdsWatcher{}
		go wc.ServeSocket(sockFile)
		wc.KeepEyesOnKeyWithPrefix(prefix)
	},
}
