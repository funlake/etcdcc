package routes

import (
	"etcdcc/apiserver/pkg/controllers"
	"github.com/astaxie/beego"
)
func init(){
	//Route = fasthttprouter.New()
	homeController := &controllers.HomeController{}
	beego.Router("/",homeController,"get:Home")
	//Route.GET("/", homeController.Home)
	//Route.POST("/find",etcdController.Get)
	//Route.POST("/set",etcdController.Put)
	//
	configController := &controllers.ConfigController{}
	//Route.POST("/config",configController.Create)
	beego.Router("/config",configController,"get:List")
	beego.Router("/config",configController,"post:Create")
	beego.Router("/config",configController,"put:Update")
	beego.Router("/config/:id",configController,"delete:Delete")
}