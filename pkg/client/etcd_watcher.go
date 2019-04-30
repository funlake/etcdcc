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
	StoreDir string
}
func (ecw *EtcdClientWatcher) KeepEyesOnKey(key string) {}

func (ecw *EtcdClientWatcher) KeepEyesOnKeyWithPrefix(key string) {
	adapter := etcd.Adapter{}
	storeDir := ecw.StoreDir + "/" + strings.Split(key,"/")[0]
	err := os.Mkdir(storeDir, 0755)
	if err != nil{
		log.Error("Can not create directory for configuration files : "+ err.Error())
		return
	}
	worker := &SyncWorker{
		storeDir: storeDir,
		shmfile:      strings.Replace(key,"/","_",-1),
		retrySeconds: ecw.RetrySeconds,
		//big queue shared by single worker
		dispatcher:   jobworker.NewBlockingDispather(1, 2000),
	}
	go worker.retryFails()
	ctx, cancel := context.WithCancel(context.Background())
	log.Info(fmt.Sprintf("Initialize configuration with %s", key))
	allKeys,err := adapter.GetMetaCacheHandler().GetStore().Get(key+"/",clientv3.WithPrefix())
	if err == nil {
		for _,e := range allKeys.(*clientv3.GetResponse).Kvs {
			sk := strings.TrimLeft(string(e.Key),key)
			ecw.ModifyLocal(sk, string(e.Value))
		}
		worker.Do(ecw.configs)
	} else {
		log.Error(err.Error())
		return
	}
	log.Info(fmt.Sprintf("Watching key with %s", key))
	for v := range adapter.GetMetaCacheHandler().GetStore().Watch(ctx, key, clientv3.WithPrefix()) {
		if v.Err() != nil {
			continue
		}
		for _, e := range v.Events {
			tp := fmt.Sprintf("%v", e.Type)
			switch tp {
			case "PUT":
				sk := strings.TrimLeft(string(e.Kv.Key),key)
				ecw.ModifyLocal(sk, string(e.Kv.Value))
			case "DELETE":
				if string(e.Kv.Key) == key {
					ecw.configs = sync.Map{}
					cancel()
				} else {
					ecw.configs.Delete(key)
				}
			}
			worker.Do(ecw.configs)
		}
	}
	cancel()
}
func (ecw *EtcdClientWatcher) ModifyLocal(key, val string) {
	v, ok := ecw.configs.Load(key)
	if ok && v.(string) != val {
		ecw.configs.Store(key, val)
		return
	} else {
		ecw.configs.Store(key, val)
	}
}
