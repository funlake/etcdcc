package utils

import (
	"github.com/astaxie/beego/orm"
	"strings"
)

//Transform dto fields to entity fields
func TransformFieldsCdt(cdt []string, fields map[string]interface{}) map[string]interface{} {
	var finals = map[string]interface{}{}
	for _, v := range cdt {
		cdv := strings.Split(v, "=")
		if f, ok := fields[cdv[0]]; ok {
			finals[f.(string)] = cdv[1]
		}
	}
	return finals
}

//Fit several conditions
func TransformQset(qs orm.QuerySeter, k string, v string) orm.QuerySeter {
	if strings.Index(v, "~") == 0 {
		qs = qs.Filter(k+"__startswith", v[1:])
	} else {
		qs = qs.Filter(k, v)
	}
	return qs
}
