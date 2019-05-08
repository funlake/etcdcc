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

func (ecw *EtcdClientWatcher) KeepEyesOnKeyWithPrefix(key string) {
	adapter := etcd.Adapter{}
	storeDir := ecw.StoreDir + "/" + key
	err := os.MkdirAll(storeDir, 0755)
	if err != nil {
		log.Error("Can not create directory for configuration files : " + err.Error())
		return
	}
	syncWorker := &SyncWorker{
		storeDir:     storeDir,
		shmfile:      strings.Replace(key, "/", "_", -1),
		retrySeconds: ecw.RetrySeconds,
		//Big queue shared with single syncWorker
		dispatcher: jobworker.NewBlockingDispather(1, 200),
	}
	go syncWorker.retryFails()
	ctx, cancel := context.WithCancel(context.Background())
	log.Info(fmt.Sprintf("Initialize configuration with %s", key))
	//Initialize all configurations under mod
	allKeys, err := adapter.GetMetaCacheHandler().GetStore().Get(key+"/", clientv3.WithPrefix())
	if err == nil {
		for _, e := range allKeys.(*clientv3.GetResponse).Kvs {
			sk := strings.TrimLeft(string(e.Key), key)
			ecw.ModifyLocal(sk, string(e.Value))
		}
		syncWorker.SyncAll(ecw.configs)
	} else {
		log.Error(err.Error())
		return
	}
	log.Info(fmt.Sprintf("Watching key with %s", key))
	//Watching mod's configurations
	for v := range adapter.GetMetaCacheHandler().GetStore().Watch(ctx, key, clientv3.WithPrefix()) {
		if v.Err() != nil {
			continue
		}
		for _, e := range v.Events {
			tp := fmt.Sprintf("%v", e.Type)
			switch tp {
			case "PUT":
				sk := strings.TrimLeft(string(e.Kv.Key), key)
				ecw.ModifyLocal(sk, string(e.Kv.Value))
			case "DELETE":
				if string(e.Kv.Key) == key {
					ecw.configs = sync.Map{}
					cancel()
				} else {
					ecw.configs.Delete(key)
				}
			}
			syncWorker.SyncAll(ecw.configs)
		}
	}
	cancel()
}
func (ecw *EtcdClientWatcher) ModifyLocal(key, val string) {
	ecw.configs.Store(key, val)
}
