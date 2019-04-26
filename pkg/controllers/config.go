package controllers

import (
	"etcdcc/apiserver/pkg/dto"
	"etcdcc/apiserver/pkg/services"
)

type ConfigController struct {
	BaseController
}
func (c *ConfigController) List(){
	service := services.Config{}
	rows,count := service.List(0,20,nil)
	c.response(RESPOK,"",map[string]interface{}{
		"result" : rows,
		"total" : count,
	})
}
func (c *ConfigController) Create()  {
	cdto := &dto.ConfigAddDto{}
	c.parseAndValidate(cdto)
	service := services.Config{}
	id,err := service.Create(cdto)
	if err != nil{
		c.fail(err.Error())
		return
	}
	c.ok(id)
}