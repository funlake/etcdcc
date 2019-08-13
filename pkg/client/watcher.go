package client

import (
	"context"
	"etcdcc/pkg/log"
	"etcdcc/pkg/storage/adapter/etcd"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/funlake/gopkg/cache"
	"strings"
)
const (
	typeYaml = "yaml"
	typeJson = "json"
	typeToml = "toml"
	typeProp = "prop"
)
// Watcher : Watching key & value change
type Watcher interface {
	KeepEyesOnKey(key string)
	KeepEyesOnKeyWithPrefix(module string)
}

//GeneralWatcher : base struct of watcher
type GeneralWatcher struct{
	//Tc *cache.TimerCacheEtcd
}

//KeepEyesOnKey : Specific key watcher
func (gw *GeneralWatcher) KeepEyesOnKey(key string) {}

//KeepEyesOnKeyWithPrefix : Specific prefix watcher
func (gw *GeneralWatcher) KeepEyesOnKeyWithPrefix(module string) {}

func (gw *GeneralWatcher) GetTimerCache() *cache.TimerCacheEtcd {
	return etcd.GetMetaCacheHandler()
}

//Init : Initialize configurations from storage while server's up
func (gw *GeneralWatcher) Init(tc *cache.TimerCacheEtcd,prefix string, callback func(k, v string)) {
	log.Info(fmt.Sprintf("Initialize configuration for prefix %s", prefix))
	//adapter := etcd.Adapter{}
	allKeys, err := tc.GetStore().Get(prefix+"/", clientv3.WithPrefix())
	if err == nil {
		for _, e := range allKeys.(*clientv3.GetResponse).Kvs {
			sk := strings.TrimPrefix(string(e.Key), prefix+"/")
			callback(sk, string(e.Value))
		}
	} else {
		log.Error(err.Error())
		return
	}
}

//Watch : Watching configuration's changes
func (gw *GeneralWatcher) Watch(tc *cache.TimerCacheEtcd,key string, putCallback func(k, v string), delCallBack func(mk, k string, cancel context.CancelFunc)) {
	//adapter := etcd.Adapter{}
	ctx, cancel := context.WithCancel(context.Background())
	log.Info(fmt.Sprintf("Watching key with %s", key))
	//Watching mod's configurations
	for v := range tc.GetStore().Watch(ctx, key, clientv3.WithPrefix()) {
		if v.Err() != nil {
			continue
		}
		for _, e := range v.Events {
			tp := fmt.Sprintf("%v", e.Type)
			sk := strings.TrimPrefix(string(e.Kv.Key), key+"/")
			switch tp {
			case "PUT":
				putCallback(sk, string(e.Kv.Value))
			case "DELETE":
				delCallBack(key, sk, cancel)
			}
		}
	}
	cancel()
}
