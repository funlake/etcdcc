package client

import (
	"context"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/funlake/gopkg/jobworker"
	"sync"
)

type EtcdClientWatcher struct {
	configs sync.Map
}

func (ecw *EtcdClientWatcher) KeepEyesOnKey(key string) {}

func (ecw *EtcdClientWatcher)KeepEyesOnKeyWithPrefix(key string,prefix interface{}) {
	adapter := etcd.Adapter{}
	worker := &SyncWorker{
		//big queue shared by single worker
		//let things
		dispatcher : jobworker.NewBlockingDispather(1,2000),
	}
	ctx,cancel := context.WithCancel(context.Background())
	log.Info(fmt.Sprintf("Watching key with %s",key))
	for v := range adapter.GetMetaCacheHandler().GetStore().Watch(ctx,key,prefix.(clientv3.OpOption)){
		if v.Err() != nil{
			continue
		}
		for _,e := range v.Events{
			tp := fmt.Sprintf("%v",e.Type)
			switch tp {
			case "PUT":
				ecw.ModifyLocal(string(e.Kv.Key),string(e.Kv.Value))
			case "DELETE" :
				if string(e.Kv.Key) == key{
					ecw.configs = sync.Map{}
					cancel()
				}else{
					ecw.configs.Delete(key)
				}
			}
			worker.Do(ecw.configs)
		}
	}
	cancel()
}

func (ecw *EtcdClientWatcher)ModifyLocal(key,val string) {
	v,ok :=  ecw.configs.Load(key)
	if ok && v.(string) != val{
		ecw.configs.Store(key,val)
		return
	}else{
		ecw.configs.Store(key,val)
	}
}