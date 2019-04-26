package mysql

import (
	"etcdcc/apiserver/pkg/log"
	"etcdcc/apiserver/pkg/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)
const MAXROWS  = 999999999
const PAGEROWS  = 20
type Adapter struct{

}
func (md Adapter) Connect(){
	err := orm.RegisterDataBase("default", "mysql", "root:-@-(-:3306)/config?charset=utf8mb4&&loc=Local")
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