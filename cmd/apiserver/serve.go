package apiserver

import (
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"etcdcc/apiserver/pkg/dao/adapter/mysql"
	"etcdcc/apiserver/pkg/log"
	_ "etcdcc/apiserver/pkg/routes"
	"github.com/astaxie/beego"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)
var (
	port string
	EtcdCertFile string
	EtcdKeyFile string
	EtcdCaFile string
	EtcdServerName string
	EtcdHosts string
	LogLevel uint8
)

var ServeCommand = &cobra.Command{
	Use: "start restful server",
	Short: "Run etcd admin restful api",
	Run: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(zerolog.Level(LogLevel))
		log.Info("Listening on port:"+port)
		etcd.Adapter.Connect(etcd.Adapter{},EtcdHosts,EtcdCertFile,EtcdKeyFile,EtcdCaFile,EtcdServerName)
		log.Info("Successfully connected to etcd server")
		mysql.Adapter.Connect(mysql.Adapter{})
		log.Info("Successfully connected to mysql server")
		beego.SetLevel(beego.LevelWarning)
		beego.Run(":"+port)
	},
}
func init()  {
	var sp = ServeCommand.PersistentFlags()
	sp.StringVar(&port,"port","80","Port of restful server")
	sp.StringVar(&EtcdCertFile,"c","/keys/ca.pem","Cert file for etcd connection")
	sp.StringVar(&EtcdKeyFile,"k","/keys/ca-key.pem","Key file for etcd connection")
	sp.StringVar(&EtcdCaFile,"ca","/keys/ca.crt","Ca file for etcd connection")
	sp.StringVar(&EtcdServerName,"sn","etcchebao","Hostname for ssl verification")
	sp.StringVar(&EtcdHosts,"hosts","127.0.0.1:2379","Hosts of etcd server")
	sp.Uint8Var(&LogLevel,"loglevel",0,"log level")
	if cobra.MarkFlagRequired(sp, "hosts") != nil {
		//cobra.MarkFlagRequired(sp, "k") != nil ||
		//cobra.MarkFlagRequired(sp, "ca") != nil {
		log.Error("Fail to set required")
	}
}

