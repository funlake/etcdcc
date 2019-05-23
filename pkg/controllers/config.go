package controllers

import (
	"etcdcc/apiserver/pkg/dto"
	"etcdcc/apiserver/pkg/services"
	"strconv"
)

//ConfigController : controller of configuration
type ConfigController struct {
	BaseController
}

//List with pagination params
func (c *ConfigController) List() {
	service := services.Config{}
	rows, count := service.List(0, 20, nil)
	c.response(RESPOK, "", map[string]interface{}{
		"result": rows,
		"total":  count,
	})
}

//Create by dto
func (c *ConfigController) Create() {
	cdto := &dto.ConfigAddDto{}
	if c.parseAndValidate(cdto) {
		service := services.Config{}
		id, err := service.Create(cdto)
		if err != nil {
			c.fail(err.Error())
			return
		}
		c.ok(id)
	}
}

//Update by dto
func (c *ConfigController) Update() {
	cdto := &dto.ConfigEditDto{}
	if c.parseAndValidate(cdto) {
		service := services.Config{}
		err := service.Update(cdto)
		if err != nil {
			c.fail(err.Error())
			return
		}
		c.ok(nil)
	}
}

//Delete by dto
func (c *ConfigController) Delete() {
	cdto := &dto.ConfigDelDto{}
	cdto.Id, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	if c.parseAndValidate(cdto) {
		service := services.Config{}
		err := service.Delete(cdto)
		if err != nil {
			c.fail(err.Error())
			return
		}
		c.ok(nil)
	}
}
