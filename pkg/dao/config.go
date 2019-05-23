package dao

import (
	"encoding/base64"
	"etcdcc/apiserver/pkg/log"
)

//CenterConfig : Table structure
type CenterConfig struct {
	Id      int    `json:"id"`
	Env     string `json:"env" orm:"column(env)"`
	Mod     string `json:"mod"`
	Key     string `json:"key"`
	Val     string `json:"val"`
	Type    string `json:"type"`
	Version string `json:"version"`
	BaseDao
}

//TableName : Return table's name in database
func (cc *CenterConfig) TableName() string {
	return "center_config"
}

//List with params of pagination & search conditions
func (cc *CenterConfig) List() (interface{}, int64) {
	db := cc.getDb()
	qs := db.QueryTable(cc.TableName())
	var rows []*CenterConfig
	var c int64
	qs = cc.filterSearch(qs, cc.q)
	_, err := qs.Limit(cc.limit, cc.start).All(&rows)
	if err == nil {
		c, _ = qs.Count()
	}
	return rows, c
}

//Find : Get records by id
func (cc *CenterConfig) Find() error {
	db := cc.getDb()
	return db.Read(cc)
}

//Create record according to dto
func (cc *CenterConfig) Create() (int64, error) {
	var insertId int64
	db := cc.getDb()
	err := db.Begin()
	insertId, err = db.Insert(cc)
	if err == nil {
		_, err = MetaCache.Put(MetaCache{}, cc.formatEtcdKeys(), base64.StdEncoding.EncodeToString([]byte(cc.Val)))
		if err != nil {
			log.Error("Etcd put error:" + err.Error())
			err = db.Rollback()
		} else {
			err = db.Commit()
		}
	} else {
		err = db.Rollback()
	}
	return insertId, err
}

//Update record according to dto
func (cc *CenterConfig) Update() error {
	db := cc.getDb()
	err := db.Begin()
	_, err = db.Update(cc)
	if err == nil {
		_, err = MetaCache.Put(MetaCache{}, cc.formatEtcdKeys(), base64.StdEncoding.EncodeToString([]byte(cc.Val)))
		if err != nil {
			log.Error("Etcd put error:" + err.Error())
			err = db.Rollback()
		} else {
			err = db.Commit()
		}
	} else {
		err = db.Rollback()
	}
	return err
}

//Delete record according to dto
func (cc *CenterConfig) Delete() error {
	db := cc.getDb()
	err := db.Begin()
	_, err = db.Delete(cc)
	if err == nil {
		_, err = MetaCache.Delete(MetaCache{}, cc.formatEtcdKeys())
		if err != nil {
			log.Error("Etcd del error:" + err.Error())
			err = db.Rollback()
		} else {
			err = db.Commit()
		}
	} else {
		err = db.Rollback()
	}
	return err
}

func (cc *CenterConfig) formatEtcdKeys() string {
	return cc.Env + "/" + cc.Mod + "/" + cc.Type + "/" + cc.Key
}
