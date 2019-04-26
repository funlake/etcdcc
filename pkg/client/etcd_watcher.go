package client

import (
	"context"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"sync"
)

type EtcdClientWatcher struct {
	configs sync.Map
}

func (ecw *EtcdClientWatcher) KeepEyesOnKey(key string) {}

func (ecw *EtcdClientWatcher)KeepEyesOnKeyWithPrefix(key string,prefix interface{}) {
	etcdAdapter := etcd.Adapter{}
	ctx,cancel := context.WithCancel(context.Background())
	for v := range etcdAdapter.GetMetaCacheHandler().GetStore().Watch(ctx,key,prefix.(clientv3.OpOption)){
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