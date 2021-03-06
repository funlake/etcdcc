package file

import (
	"github.com/funlake/etcdcc/pkg/client"
	"github.com/funlake/etcdcc/pkg/log"
	"github.com/funlake/etcdcc/pkg/storage/adapter/etcd"
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
	storeDir       string
)

func init() {
	var cp = FileCommand.PersistentFlags()
	cp.StringVar(&prefix, "prefix", "global", "Name of prefix of current module")
	cp.StringVar(&etcdCertFile, "c", "/keys/client.pem", "Cert file for etcd connection")
	cp.StringVar(&etcdKeyFile, "k", "/keys/client-key.pem", "Key file for etcd connection")
	cp.StringVar(&etcdCaFile, "ca", "/keys/ca.pem", "Ca file for etcd connection")
	cp.StringVar(&etcdServerName, "sn", "", "ServerName for ssl verification")
	cp.StringVar(&etcdHosts, "hosts", "127.0.0.1:2379", "Hosts of etcd server")
	cp.IntVar(&retrySeconds, "retrySeconds", 3, "Fails retry in ? seconds")
	cp.StringVar(&storeDir, "storeDir", "/tmp", "Directory of config file")

	if cobra.MarkFlagRequired(cp, "prefix") != nil ||
		cobra.MarkFlagRequired(cp, "hosts") != nil {
		//cobra.MarkFlagRequired(sp, "k") != nil ||
		//cobra.MarkFlagRequired(sp, "ca") != nil {
		log.Error("Fail to set required")
	}
}

//FileCommand : file storage watching command
var FileCommand = &cobra.Command{
	Use:   "client.file",
	Short: "Listening config changes & modified local configuration",
	Run: func(cmd *cobra.Command, args []string) {
		etcd.Connect(etcdHosts, etcdCertFile, etcdKeyFile, etcdCaFile, etcdServerName)
		log.Info("Successfully connected to etcd server")
		wc := &client.EtcdFileWatcher{RetrySeconds: retrySeconds, StoreDir: storeDir, Tc: etcd.GetMetaCacheHandler()}
		wc.KeepEyesOnKeyWithPrefix(prefix)
	},
}
