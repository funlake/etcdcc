package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strings"
)

const (
	RESPOK = iota
	RESPFAIL
)

//type CommonResponse struct {
//	Code int    `json:"code"`
//	Msg  string `json:"msg"`
//	Data interface{} `json:"data"`
//}

type BaseController struct {
	beego.Controller
}

func (b *BaseController) response(code int, msg string, data ...interface{}) {
	//ctx.Response.Header.Add("Content-Type","application/json; charset=utf-8")
	//res,_ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&CommonResponse{
	//	RESPONSE_OK,msg,data,
	//})
	//_,err := ctx.Write(res)
	//if err != nil{
	//	log.Error(err.Error())
	//}
	resp := make(map[string]interface{})
	resp["code"] = code
	resp["msg"] = msg
	if len(data) >= 1 {
		resp["data"] = data[0]
	}
	if len(data) >= 2 {
		resp["total"] = data[1]
	}
	b.Data["json"] = resp
	b.ServeJSON()
}

func (b *BaseController) ok(data ...interface{}) {
	b.response(RESPOK, "ok", data[0])
}

func (b *BaseController) fail(msg string) {
	resp := make(map[string]interface{})
	resp["code"] = RESPFAIL
	resp["msg"] = msg
	b.Data["json"] = resp
	b.ServeJSON()
}

func (b *BaseController) parseAndValidate(obj interface{}) bool {
	if err := b.ParseForm(obj); err != nil {
		b.fail(err.Error())
		return false
	}
	valid := &validation.Validation{}
	if v, _ := valid.Valid(obj); !v {
		var errs []string
		for _, err := range valid.Errors {
			errs = append(errs, strings.ToLower(strings.Split(err.Key, ".")[0])+":"+err.Message)
		}
		b.fail(strings.Join(errs, ","))
		return false
	}
	return true
}
