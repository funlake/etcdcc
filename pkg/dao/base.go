package dao

import (
	"etcdcc/apiserver/pkg/utils"
	"github.com/astaxie/beego/orm"
	"time"
)

//Max amount of rows returns per page
const MAXROWS = 999999999

//Normal amount of rows return per page
const PAGEROWS = 20

//BaseDao : Wrap of dao
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
	return orm.NewOrm()
}

//Extends q for more search conditions
func (bd *BaseDao) filterSearch(qs orm.QuerySeter, q []string) orm.QuerySeter {
	if len(q) > 0 {
		for k, v := range utils.TransformFieldsCdt(q, bd.searchMap) {
			qs = utils.TransformQset(qs, k, v.(string))
		}
	}
	return qs
}

//SetSearchMap : Set search conditions
func (bd *BaseDao) SetSearchMap(sm map[string]interface{}) {
	bd.searchMap = sm
}

//SetPageParams : Set pagination params
func (bd *BaseDao) SetPageParams(start int, limit int) {
	bd.start = start
	bd.limit = limit
}

//SetSearchCdt : Set search conditions
func (bd *BaseDao) SetSearchCdt(q []string) {
	bd.q = q
}
