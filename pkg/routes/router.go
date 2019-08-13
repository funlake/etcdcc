package routes

import (
	"etcdcc/pkg/controllers"
	"github.com/astaxie/beego"
)

func init() {
	homeController := &controllers.HomeController{}
	beego.Router("/", homeController, "get:Home")
	configController := &controllers.ConfigController{}
	beego.Router("/configs", configController, "get:List")
	beego.Router("/config", configController, "post:Create")
	beego.Router("/config/:id", configController, "put:Update")
	beego.Router("/config/:id", configController, "delete:Delete")
}
