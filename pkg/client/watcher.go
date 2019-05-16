package client

import (
	"context"
	"etcdcc/apiserver/pkg/storage/adapter/etcd"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

// @Summary
// 1. 同步器用job worker异步阻塞模式，只开一个工作协程，避免写锁,协程panic重启由类库自动实现
// 2. 定时重试
// 3. 只兼容linux,一些命令用linux原生执行(通用后期再考虑)
type Watcher interface {
	KeepEyesOnKey(key string)
	KeepEyesOnKeyWithPrefix(module string)
	//Init(moduleKey string, callback func(k, v string))
	//Watch(key string, putCallback func(k, v string), delCallBack func(mk, k string, cancel context.CancelFunc))
}

type GeneralWatcher struct{}

func (gw *GeneralWatcher) KeepEyesOnKey(key string)              {}
func (gw *GeneralWatcher) KeepEyesOnKeyWithPrefix(module string) {}
func (gw *GeneralWatcher) Init(prefix string, callback func(k, v string)) {
	log.Info(fmt.Sprintf("Initialize configuration for prefix %s", prefix))
	adapter := etcd.Adapter{}
	allKeys, err := adapter.GetMetaCacheHandler().GetStore().Get(prefix+"/", clientv3.WithPrefix())
	if err == nil {
		for _, e := range allKeys.(*clientv3.GetResponse).Kvs {
			sk := strings.TrimPrefix(string(e.Key), prefix+"/")
			//ecw.ModifyLocal(sk, string(e.Value))
			callback(sk, string(e.Value))
		}
		//syncWorker.SyncAll(ecw.configs)
	} else {
		log.Error(err.Error())
		return
	}
}
func (gw *GeneralWatcher) Watch(key string, putCallback func(k, v string), delCallBack func(mk, k string, cancel context.CancelFunc)) {
	adapter := etcd.Adapter{}
	ctx, cancel := context.WithCancel(context.Background())
	log.Info(fmt.Sprintf("Watching key with %s", key))
	//Watching mod's configurations
	for v := range adapter.GetMetaCacheHandler().GetStore().Watch(ctx, key, clientv3.WithPrefix()) {
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
