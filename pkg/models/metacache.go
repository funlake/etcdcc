package models

import (
	"errors"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
)

var metaHandler = etcd.Adapter{}
type MetaCache struct {}
func (e MetaCache) Put(key,val string)(interface{},error){
	return metaHandler.GetMetaCacheHandler().GetStore().Set(key,val)
}
func (e MetaCache) Delete(key string)(interface{},error){
	return metaHandler.GetMetaCacheHandler().GetStore().Delete(key)
}
func (e MetaCache) Get(key string)(interface{},error){
	//Automatically add a local cache + watch
	val ,err := metaHandler.GetMetaCacheHandler().Get("",key,0)
	if err == nil && val == ""{
		err = errors.New("Value Not set")
	}
	return val,err
}