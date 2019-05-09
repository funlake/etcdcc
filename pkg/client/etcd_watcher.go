package client

import (
	"context"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/funlake/gopkg/jobworker"
	"os"
	"strings"
	"sync"
)

type EtcdClientWatcher struct {
	configs      sync.Map
	RetrySeconds int
	StoreDir     string
}

func (ecw *EtcdClientWatcher) KeepEyesOnKey(key string) {}

func (ecw *EtcdClientWatcher) KeepEyesOnKeyWithPrefix(moduleKey string) {
	storeDir := ecw.StoreDir + "/" + moduleKey
	err := os.MkdirAll(storeDir, 0755)
	if err != nil {
		log.Error("Can not create directory for configuration files : " + err.Error())
		return
	}
	syncWorker := &SyncWorker{
		storeDir:     storeDir,
		shmfile:      strings.Replace(moduleKey, "/", "_", -1),
		retrySeconds: ecw.RetrySeconds,
		//Big queue shared with single syncWorker
		dispatcher: jobworker.NewBlockingDispather(1, 200),
	}
	go syncWorker.retryFails()
	//Initialize all configurations under mod
	ecw.Init(moduleKey, func(k, v string) {
		syncWorker.SyncOne(k,v)
	})
	ecw.Watch(moduleKey, func(k, v string) {
		syncWorker.SyncOne(k,v)
	}, func(mk,k string,cancel context.CancelFunc) {
		//监听key == 删除key，则整个watch停止
		if mk == k {
			cancel()
		}
		syncWorker.RemoveOne(k)
	})
}
func (ecw *EtcdClientWatcher) Init(moduleKey string,callback func(k,v string)){
	log.Info(fmt.Sprintf("Initialize configuration with %s", moduleKey))
	adapter := etcd.Adapter{}
	allKeys, err := adapter.GetMetaCacheHandler().GetStore().Get(moduleKey+"/", clientv3.WithPrefix())
	if err == nil {
		for _, e := range allKeys.(*clientv3.GetResponse).Kvs {
			sk := strings.TrimPrefix(string(e.Key), moduleKey+"/")
			ecw.ModifyLocal(sk, string(e.Value))
			callback(sk,string(e.Value))
		}
		//syncWorker.SyncAll(ecw.configs)
	} else {
		log.Error(err.Error())
		return
	}
}
func (ecw *EtcdClientWatcher) Watch(key string,putCallback func(k,v string),delCallBack func(mk,k string,cancel context.CancelFunc)){
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
				putCallback(sk,string(e.Kv.Value))
			case "DELETE":
				delCallBack(key,sk,cancel)
			}
		}
	}
	cancel()
}
func (ecw *EtcdClientWatcher) ModifyLocal(key, val string) {
	ecw.configs.Store(key, val)
}
