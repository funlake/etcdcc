package dto

import "github.com/astaxie/beego/validation"

var CONFIG_SEARCH = map[string]interface{}{
	"e": "env",
	"k": "key",
	"m": "mod",
}

// @规范 : DTO 的字段必须与form标识的标签一一对应，允许大小写关系，但不能不一致，否则验证信息会让人疑惑
type ConfigAddDto struct {
	Env  string `form:"env" valid:"Required"`
	Mod  string `form:"mod" valid:"Required"`
	Key  string `form:"key" valid:"Required"`
	Val  string `form:"val" valid:"Required"`
	Type string `form:"type"`
}

// 自定义额外的判定
func (cc *ConfigAddDto) Valid(v *validation.Validation) {
	//if strings.Index(u.Name, "admin") != -1 {
	//	// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
	//	v.SetError("Name", "名称里不能含有 admin")
	//}
}

type ConfigEditDto struct {
	Id   int    `form:"id" valid:"Required;Min(1)"`
	Env  string `form:"env" valid:"Required"`
	Mod  string `form:"mod" valid:"Required"`
	Key  string `form:"key" valid:"Required"`
	Val  string `form:"val" valid:"Required"`
	Type string `form:"type"`
}
type ConfigDelDto struct {
	Id int `form:"id" valid:"Required;Min(1)"`
}
