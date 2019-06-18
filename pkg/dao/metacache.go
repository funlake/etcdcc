package dao

import (
	"errors"
	"etcdcc/pkg/storage/adapter/etcd"
)

var metaHandler = etcd.Adapter{}

//MetaCache : Meta data storage
type MetaCache struct{}

//Put : Update key with value
func (e MetaCache) Put(key, val string) (interface{}, error) {
	return metaHandler.GetMetaCacheHandler().GetStore().Set(key, val)
}

//Delete by key
func (e MetaCache) Delete(key string) (interface{}, error) {
	return metaHandler.GetMetaCacheHandler().GetStore().Delete(key)
}

//Get specific key
func (e MetaCache) Get(key string) (interface{}, error) {
	val, err := metaHandler.GetMetaCacheHandler().Get("", key, 0)
	if err == nil && val == "" {
		err = errors.New("Value not set")
	}
	return val, err
}
