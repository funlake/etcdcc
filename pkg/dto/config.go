package dto

import "github.com/astaxie/beego/validation"

//CONFIG_SEARCH : search condition of config list
var CONFIG_SEARCH = map[string]interface{}{
	"e": "env",
	"k": "key",
	"m": "mod",
}

// ConfigAddDto : struct property's name should identical with  value of form ,or warning information may confuse people
type ConfigAddDto struct {
	Env  string `form:"env" valid:"Required"`
	Mod  string `form:"mod" valid:"Required"`
	Key  string `form:"key" valid:"Required"`
	Val  string `form:"val" valid:"Required"`
	Type string `form:"type"`
}

// Valid : extra validate
func (cc *ConfigAddDto) Valid(v *validation.Validation) {
	//if strings.Index(u.Name, "admin") != -1 {
	//	// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
	//	v.SetError("Name", "名称里不能含有 admin")
	//}
}

//ConfigEditDto : edit dto
type ConfigEditDto struct {
	Id   int    `form:"id" valid:"Required;Min(1)"`
	Env  string `form:"env" valid:"Required"`
	Mod  string `form:"mod" valid:"Required"`
	Key  string `form:"key" valid:"Required"`
	Val  string `form:"val" valid:"Required"`
	Type string `form:"type"`
}

//ConfigDelDto : delete dto
type ConfigDelDto struct {
	Id int `form:"id" valid:"Required;Min(1)"`
}
