package mysql

import (
	"etcdcc/apiserver/pkg/log"
	"etcdcc/apiserver/pkg/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)
type Adapter struct{

}
func (md Adapter) Connect(){
	viper.AutomaticEnv()
	host := viper.Get("MYSQL_HOST")
	user := viper.Get("MYSQL_USERNAME")
	pwd  := viper.Get("MYSQL_PASSWORD")
	db   := viper.Get("MYSQL_DB_CC")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&&loc=Local",user,pwd,host,db)
	err := orm.RegisterDataBase("default", "mysql", dsn)
	if err != nil {
		log.Fatal("Mysql connection fail:"+err.Error())
	}else{
		//设置最大链接数
		orm.SetMaxIdleConns("default", 100)
		orm.SetMaxOpenConns("default", 300)
		orm.RegisterModel(&models.CenterConfig{})
		orm.Debug = true
	}
}