package server

import (
	"github.com/astaxie/beego"
	"github.com/funlake/etcdcc/pkg/log"
	//beego need import routes here
	_ "github.com/funlake/etcdcc/pkg/routes"
	"github.com/funlake/etcdcc/pkg/storage/adapter/etcd"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	port           string
	etcdCertFile   string
	etcdKeyFile    string
	etcdCaFile     string
	etcdServerName string
	etcdHosts      string
	logLevel       uint8
)

//ServeCommand : server command for configurations curd
var ServeCommand = &cobra.Command{
	Use:   "server.start",
	Short: "Run etcd admin restful api",
	Run: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.Level(logLevel))
		log.Info("Listening on port:" + port)
		etcd.Connect(etcdHosts, etcdCertFile, etcdKeyFile, etcdCaFile, etcdServerName)
		log.Info("Successfully connected to etcd server")
		//mysql.Adapter.Connect(mysql.Adapter{})
		//log.Info("Successfully connected to mysql server")
		beego.SetLevel(beego.LevelWarning)
		beego.Run(":" + port)
	},
}

func init() {
	var sp = ServeCommand.PersistentFlags()
	sp.StringVar(&port, "port", "80", "Port of restful server")
	sp.StringVar(&etcdCertFile, "c", "/keys/client.pem", "Cert file for etcd connection")
	sp.StringVar(&etcdKeyFile, "k", "/keys/client-key.pem", "Key file for etcd connection")
	sp.StringVar(&etcdCaFile, "ca", "/keys/ca.pem", "Ca file for etcd connection")
	sp.StringVar(&etcdServerName, "sn", "", "ServerName for ssl verification")
	sp.StringVar(&etcdHosts, "hosts", "127.0.0.1:2379", "Hosts of etcd server")
	sp.Uint8Var(&logLevel, "loglevel", 0, "Log level")
	if cobra.MarkFlagRequired(sp, "hosts") != nil {
		//cobra.MarkFlagRequired(sp, "k") != nil ||
		//cobra.MarkFlagRequired(sp, "ca") != nil {
		log.Error("Fail to set required")
	}
}
