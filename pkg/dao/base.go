package dao

import (
	"etcdcc/apiserver/pkg/utils"
	"github.com/astaxie/beego/orm"
	"time"
)

const MAXROWS = 999999999
const PAGEROWS = 20

type BaseDao struct {
	CreatedTime time.Time `orm:"auto_now_add;type(datetime)" json:"created_time"`
	UpdatedTime time.Time `orm:"auto_now;type(datetime)" json:"updated_time"`
	start       int
	limit       int
	searchMap   map[string]interface{}
	q           []string
}

//Does orm of beego get connection from connection pool every time while executing the sql?
//if not then we need to do orm.NewOrm every time.
//@found
// Seems we should not use the same orm object when doing tnx
// let's keep doing orm.NewOrm action every time
// author said it would not down the performance, see @link
//@link
// https://github.com/astaxie/beego/issues/1524
func (bd *BaseDao) getDb() orm.Ormer {
	//return onceOrm
	return orm.NewOrm()
}

//func (bd *BaseDao) fetchRows(qs orm.QuerySeter) ([]orm.Params,int64) {
//	var c int64
//	var rows []orm.Params
//	qs = bd.filterSearch(qs,bd.q,bd.searchMap)
//	_, err := qs.Limit(bd.limit, bd.start).Values(&rows)
//	if err == nil {
//		c, _ = qs.Count()
//	}
//	return rows,c
//}

func (bd *BaseDao) filterSearch(qs orm.QuerySeter, q []string) orm.QuerySeter {
	//后期加入搜索条件可利用q参数
	if len(q) > 0 {
		for k, v := range utils.TransformFieldsCdt(q, bd.searchMap) {
			qs = utils.TransformQset(qs, k, v.(string))
		}
	}
	return qs
}

//设置搜索条件与数据表字段自建的关联关系
func (bd *BaseDao) SetSearchMap(sm map[string]interface{}) {
	bd.searchMap = sm
}

//设置分页参数
func (bd *BaseDao) SetPageParams(start int, limit int) {
	bd.start = start
	bd.limit = limit
}

//设置搜索条件
func (bd *BaseDao) SetSearchCdt(q []string) {
	bd.q = q
}
